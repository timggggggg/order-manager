package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type TransactionManager interface {
	GetQueryEngine(ctx context.Context) QueryEngine
	RunReadUncommitted(ctx context.Context, fn func(ctxTx context.Context) error) error
	RunReadCommitted(ctx context.Context, fn func(ctxTx context.Context) error) error
	RunSerializable(ctx context.Context, fn func(ctxTx context.Context) error) error
}

type PgRepository struct {
	txManager TransactionManager
}

func NewPgRepository(txManager TransactionManager) *PgRepository {
	return &PgRepository{
		txManager: txManager,
	}
}

func (r *PgRepository) GetByID(ctx context.Context, id int64) (*OrderDB, error) {
	var order OrderDB

	query := `
		SELECT *
		FROM orders
		WHERE id = $1
		FOR UPDATE
	`
	tx := r.txManager.GetQueryEngine(ctx)
	err := pgxscan.Get(ctx, tx, &order, query, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, models.ErrorOrderNotFound
		}
		return nil, err
	}

	return &order, nil
}

func (r *PgRepository) GetByIDs(ctx context.Context, ids []int64) (OrdersDBMapStorage, error) {
	var orders OrdersDBSliceStorage

	query := `
		SELECT *
		FROM orders
		WHERE id = ANY($1)
		FOR UPDATE
	`
	tx := r.txManager.GetQueryEngine(ctx)
	err := pgxscan.Select(ctx, tx, &orders, query, ids)
	if err != nil {
		return nil, err
	}

	ordersMap := OrdersDBStorageSliceToMap(orders)
	for _, id := range ids {
		if _, exists := ordersMap[id]; !exists {
			return nil, fmt.Errorf("orderID=%d: %w", id, models.ErrorOrderNotFound)
		}
	}

	return ordersMap, nil
}

func (r *PgRepository) GetByUserIDCursorPagination(ctx context.Context, userID int64, limit int64, cursorID int64) (OrdersDBSliceStorage, error) {
	var orders OrdersDBSliceStorage

	query := `
		SELECT *
		FROM orders
		WHERE user_id = $1
		AND id > $2
		ORDER BY id ASC
		LIMIT $3
	`
	tx := r.txManager.GetQueryEngine(ctx)
	err := pgxscan.Select(ctx, tx, &orders, query, userID, cursorID, limit)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *PgRepository) GetAll(ctx context.Context, limit int64, offset int64) (OrdersDBSliceStorage, error) {
	var orders OrdersDBSliceStorage

	query := `
		SELECT *
		FROM orders
		ORDER BY id ASC
		LIMIT $1 OFFSET $2
	`
	tx := r.txManager.GetQueryEngine(ctx)
	err := pgxscan.Select(ctx, tx, &orders, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *PgRepository) GetReturnsLimitOffsetPagination(ctx context.Context, limit int64, offset int64) (OrdersDBSliceStorage, error) {
	var orders OrdersDBSliceStorage

	query := `
		SELECT *
		FROM orders
		WHERE order_status = 'returned'
		ORDER BY id ASC
		LIMIT $1 OFFSET $2
	`
	tx := r.txManager.GetQueryEngine(ctx)
	err := pgxscan.Select(ctx, tx, &orders, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *PgRepository) CreateOrder(ctx context.Context, order *OrderDB) (err error) {
	query := `
		INSERT INTO orders (id, user_id,
		order_status, accept_time, expire_time, issue_time,
		weight, cost, package, extra_package)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	tx := r.txManager.GetQueryEngine(ctx)

	_, err = tx.Exec(ctx, query,
		order.ID, order.UserID, order.Status, order.AcceptTime, order.ExpireTime,
		order.IssueTime, order.Weight, order.Cost, order.Package, order.ExtraPackage)

	if err != nil && strings.Contains(err.Error(), "duplicate") {
		return models.ErrorOrderAlreadyExists
	}

	return err
}

func (r *PgRepository) ReturnOrder(ctx context.Context, id int64) (*OrderDB, error) {
	var order OrderDB
	query := `
		UPDATE orders SET order_status = 'returned'
		WHERE id = $1
		RETURNING *
	`
	tx := r.txManager.GetQueryEngine(ctx)
	err := pgxscan.Get(ctx, tx, &order, query, id)

	return &order, err
}

func (r *PgRepository) IssueOrders(ctx context.Context, ids []int64) (OrdersDBSliceStorage, error) {
	var orders OrdersDBSliceStorage

	query := `
		UPDATE orders
		SET issue_time = $1
		WHERE id = ANY($2)
		RETURNING *
	`
	tx := r.txManager.GetQueryEngine(ctx)
	err := pgxscan.Select(ctx, tx, &orders, query, sql.NullTime{Time: time.Now(), Valid: true}, ids)

	return orders, err
}

func (r *PgRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM orders
	  	WHERE id = $1
	`

	tx := r.txManager.GetQueryEngine(ctx)
	_, err := tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgRepository) LogOrderEvent(ctx context.Context, order *OrderDB, status models.OrderStatus) error {
	query := `
		INSERT INTO order_events (order_id, user_id, event_type)
		values ($1, $2, $3)
	`
	tx := r.txManager.GetQueryEngine(ctx)
	_, err := tx.Exec(ctx, query, order.ID, order.UserID, status)

	return err
}
