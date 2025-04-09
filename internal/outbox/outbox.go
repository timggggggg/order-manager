// nolint
package outbox

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type IProducer interface {
	SendMessage(topic string, payload []byte) error
	Close() error
}

type Outbox struct {
	pool        *pgxpool.Pool
	tableName   string
	producer    IProducer
	topic       string
	MaxAttempts int64
}

func NewOutbox(pool *pgxpool.Pool, tableName string, producer IProducer, topic string, maxAttempts int64) *Outbox {
	return &Outbox{
		pool:        pool,
		tableName:   tableName,
		producer:    producer,
		topic:       topic,
		MaxAttempts: maxAttempts,
	}
}

func (o *Outbox) CloseProducer() {
	o.producer.Close()
}

func (o *Outbox) RenewTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id is not provided", http.StatusBadRequest)
		return
	}

	ID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf(`
		UPDATE %s
		SET status = $1, updated_at = NOW(), attempts_left = $2
		WHERE id = $3
	`, o.tableName)

	_, err = o.pool.Exec(ctx,
		query,
		models.TaskStatusCreated,
		o.MaxAttempts,
		ID,
	)

	if err != nil {
		http.Error(w, fmt.Sprintf("task renew failed: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "task id=%d renewed!\n", ID)
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
		var task models.OutboxTask

		if err := rows.Scan(&task.ID, &task.Payload, &task.Status, &task.AttemptsLeft); err != nil {
			log.Printf("Scan failed: %v", err)
			continue
		}

		query = fmt.Sprintf(`
            UPDATE %s
            SET status = $3, updated_at = NOW(), attempts_left = $1
            WHERE id = $2
        `, o.tableName)

		task.AttemptsLeft -= 1

		_, err = o.pool.Exec(ctx,
			query,
			task.AttemptsLeft,
			task.ID,
			models.TaskStatusProcessing,
		)

		if err != nil {
			log.Printf("Update failed: %v", err)
			continue
		}

		err = o.producer.SendMessage(o.topic, task.Payload)

		task.Status = models.TaskStatusCompleted
		if err != nil {
			log.Printf("Kafka send failed: %v", err)
			task.Status = models.TaskStatusFailed
			if task.AttemptsLeft == 0 {
				task.Status = models.TaskStatusNoAttemptsLeft
			}
		}

		query = fmt.Sprintf(`
            UPDATE %s
            SET status = $2, updated_at = NOW()
            WHERE id = $1
        `, o.tableName)

		_, err = o.pool.Exec(ctx, query, task.ID, task.Status)
		if err != nil {
			log.Printf("Completion failed: %v", err)
		}
	}
}
