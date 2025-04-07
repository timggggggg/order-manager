package logger

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type DBLogger struct {
	pool *pgxpool.Pool
}

func NewDBLogger(pool *pgxpool.Pool) *DBLogger {
	return &DBLogger{
		pool: pool,
	}
}

func (l *DBLogger) LogStatusChange(ctx context.Context, ts time.Time, id int64, statusFrom, statusTo models.OrderStatus) {
	query := `
		INSERT INTO logs (order_id, status_from, status_to, ts)
		VALUES ($1, $2, $3, $4)
	`
	tx := l.pool
	_, err := tx.Exec(ctx, query, id, statusFrom, statusTo, ts)

	if err != nil {
		return
	}

	payload := &models.StatusChangeLog{
		OrderID:    id,
		StatusFrom: statusFrom,
		StatusTo:   statusTo,
		Timestamp:  ts,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return
	}

	task := models.NewOutboxTask(models.TaskTypeStatusChangeLog, payloadBytes, 3)

	err = outboxInsert(ctx, tx, task)
	if err != nil {
		return
	}
}

func (l *DBLogger) LogRequest(ctx context.Context, ts time.Time, method, url, request_body string) {
	query := `
		INSERT INTO http_req_logs (ts, method, url, request_body)
		VALUES ($1, $2, $3, $4)
	`
	tx := l.pool
	_, err := tx.Exec(ctx, query, ts, method, url, request_body)

	if err != nil {
		return
	}

	payload := &models.HttpReqLog{
		Timestamp:   ts,
		Method:      method,
		URL:         url,
		RequestBody: request_body,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return
	}

	task := models.NewOutboxTask(models.TaskTypeHttpReqLog, payloadBytes, 3)

	err = outboxInsert(ctx, tx, task)
	if err != nil {
		return
	}
}

func (l *DBLogger) LogResponse(ctx context.Context, ts time.Time, code int64, body string) {
	query := `
		INSERT INTO http_resp_logs (ts, status_code, body)
		VALUES ($1, $2, $3)
	`
	tx := l.pool
	_, err := tx.Exec(ctx, query, ts, code, body)

	if err != nil {
		return
	}

	payload := &models.HttpRespLog{
		Timestamp:  ts,
		StatusCode: code,
		Body:       body,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return
	}

	task := models.NewOutboxTask(models.TaskTypeHttpRespLog, payloadBytes, 3)

	err = outboxInsert(ctx, tx, task)
	if err != nil {
		return
	}
}

func outboxInsert(ctx context.Context, tx *pgxpool.Pool, task *models.OutboxTask) error {
	query := `
		INSERT INTO tasks (task_type, payload, status, attempts_left)
		VALUES ($1, $2, $3, $4)
	`

	_, err := tx.Exec(ctx, query, task.Type, task.Payload, task.Status, task.AttemptsLeft)

	return err
}
