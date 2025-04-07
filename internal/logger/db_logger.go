package logger

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type Task struct {
	OrderID     int64              `json:"order_id,omitempty"`
	StatusFrom  models.OrderStatus `json:"status_from,omitempty"`
	StatusTo    models.OrderStatus `json:"status_to,omitempty"`
	Timestamp   time.Time          `json:"ts"`
	Method      string             `json:"method,omitempty"`
	URL         string             `json:"url,omitempty"`
	RequestBody string             `json:"request_body,omitempty"`
	StatusCode  int64              `json:"status_code,omitempty"`
	Body        string             `json:"body,omitempty"`
}

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

	task := &Task{
		OrderID:    id,
		StatusFrom: statusFrom,
		StatusTo:   statusTo,
		Timestamp:  ts,
	}

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

	task := &Task{
		Timestamp:   ts,
		Method:      method,
		URL:         url,
		RequestBody: request_body,
	}

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

	task := &Task{
		Timestamp:  ts,
		StatusCode: code,
		Body:       body,
	}

	err = outboxInsert(ctx, tx, task)
	if err != nil {
		return
	}
}

func outboxInsert(ctx context.Context, tx *pgxpool.Pool, jsonData *Task) error {
	result, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO tasks (payload)
		VALUES ($1)
	`

	_, err = tx.Exec(ctx, query, result)

	return err
}
