package kafka

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxWorkerPool struct {
	pool         *pgxpool.Pool
	producer     sarama.SyncProducer
	topic        string
	interval     time.Duration
	maxRetries   int
	workersCount int
	wg           sync.WaitGroup
}

func NewOutboxWorkerPool(workersCount int, pool *pgxpool.Pool, brokers []string, topic string, interval time.Duration) (*OutboxWorkerPool, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("error creating producer: %v", err)
	}

	return &OutboxWorkerPool{
		pool:         pool,
		producer:     producer,
		topic:        topic,
		interval:     interval,
		maxRetries:   3,
		workersCount: workersCount,
	}, nil
}

func (w *OutboxWorkerPool) Start(ctx context.Context) {
	for range w.workersCount {
		w.wg.Add(1)
		go w.run(ctx)
	}
}

func (w *OutboxWorkerPool) Shutdown() {
	w.producer.Close()
	w.wg.Wait()
}

func (w *OutboxWorkerPool) run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)

	defer w.wg.Done()
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.processBatch(ctx)
		case <-ctx.Done():
			log.Println("Stopping outbox worker")
			return
		}
	}
}

func (w *OutboxWorkerPool) processBatch(ctx context.Context) {
	// tx, err := w.pool.Begin(ctx)
	// if err != nil {
	// 	log.Printf("Failed to begin transaction: %v", err)
	// 	return
	// }
	// defer tx.Rollback(ctx)

	query := `
        SELECT id, audit_log, attempts_left
        FROM outbox
        WHERE status IN ('CREATED', 'FAILED')
        AND (updated_at IS NULL OR updated_at < $1)
        AND attempts_left > 0
        FOR UPDATE SKIP LOCKED
    `

	rows, err := w.pool.Query(ctx,
		query,
		time.Now().Add(-2*time.Second),
	)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var auditJSON []byte
		var attempts int

		if err := rows.Scan(&id, &auditJSON, &attempts); err != nil {
			log.Printf("Scan failed: %v", err)
			continue
		}

		query = `
            UPDATE outbox
            SET status = 'PROCESSING', updated_at = NOW(), attempts_left = $1
            WHERE id = $2
        `

		_, err = w.pool.Exec(ctx,
			query,
			attempts-1,
			id,
		)

		if err != nil {
			log.Printf("Update failed: %v", err)
			continue
		}

		// if err := tx.Commit(ctx); err != nil {
		// 	log.Printf("Commit failed: %v", err)
		// 	return
		// }

		_, _, err = w.producer.SendMessage(&sarama.ProducerMessage{
			Topic: w.topic,
			Value: sarama.ByteEncoder(auditJSON),
		})

		if err != nil {
			log.Printf("Kafka send failed: %v", err)
			w.handleFailure(ctx, id, attempts-1)
			continue
		}

		query = `
            UPDATE outbox
            SET status = 'COMPLETED'
            WHERE id = $1
        `
		_, err = w.pool.Exec(ctx, query, id)
		if err != nil {
			log.Printf("Completion failed: %v", err)
		}
	}
}

func (w *OutboxWorkerPool) handleFailure(ctx context.Context, id int64, attemptsLeft int) {
	status := "FAILED"
	if attemptsLeft <= 0 {
		status = "NO_ATTEMPTS_LEFT"
	}

	query := ` 
        UPDATE outbox
        SET status = $1, updated_at = NOW()
        WHERE id = $2
    `
	_, err := w.pool.Exec(ctx, query, status, id)
	if err != nil {
		log.Printf("Failure update failed: %v", err)
	}
}
