package logger

import (
	"context"
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
}
