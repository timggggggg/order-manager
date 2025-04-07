package outbox

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type IProducer interface {
	SendMessage(topic string, payload []byte) error
	Close() error
}

type Outbox struct {
	pool      *pgxpool.Pool
	tableName string
	producer  IProducer
	topic     string
}

func NewOutbox(pool *pgxpool.Pool, tableName string, producer IProducer, topic string) *Outbox {
	return &Outbox{
		pool:      pool,
		tableName: tableName,
		producer:  producer,
		topic:     topic,
	}
}

func (o *Outbox) CloseProducer() {
	o.producer.Close()
}

func (o *Outbox) ProcessBatch(ctx context.Context) {
	query := fmt.Sprintf(`
        SELECT id, payload, status, attempts_left
        FROM %s
        WHERE status IN ($2, $3)
        AND (updated_at IS NULL OR updated_at < $1)
        AND attempts_left > 0
        FOR UPDATE SKIP LOCKED
    `, o.tableName)

	rows, err := o.pool.Query(ctx,
		query,
		time.Now().Add(-2*time.Second),
		models.TaskStatusCreated,
		models.TaskStatusFailed,
	)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var auditJSON []byte
		var status models.TaskStatus
		var attempts int

		if err := rows.Scan(&id, &auditJSON, &status, &attempts); err != nil {
			log.Printf("Scan failed: %v", err)
			continue
		}

		query = fmt.Sprintf(`
            UPDATE %s
            SET status = $3, updated_at = NOW(), attempts_left = $1
            WHERE id = $2
        `, o.tableName)

		attempts -= 1

		_, err = o.pool.Exec(ctx,
			query,
			attempts,
			id,
			models.TaskStatusProcessing,
		)

		if err != nil {
			log.Printf("Update failed: %v", err)
			continue
		}

		err = o.producer.SendMessage(o.topic, auditJSON)

		status = models.TaskStatusCompleted
		if err != nil {
			log.Printf("Kafka send failed: %v", err)
			status = models.TaskStatusFailed
			if attempts == 0 {
				status = models.TaskStatusNoAttemptsLeft
			}
		}

		query = fmt.Sprintf(`
            UPDATE %s
            SET status = $2, updated_at = NOW()
            WHERE id = $1
        `, o.tableName)

		_, err = o.pool.Exec(ctx, query, id, status)
		if err != nil {
			log.Printf("Completion failed: %v", err)
		}
	}
}
