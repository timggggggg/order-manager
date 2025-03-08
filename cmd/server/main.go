package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.ozon.dev/timofey15g/homework/internal/service"
	"gitlab.ozon.dev/timofey15g/homework/internal/storage/postgres"
)

func main() {
	ctx := context.Background()

	pool, err := newPgxPool(ctx)
	if err != nil {
		log.Fatal("error newPgxPool", err)
	}

	defer func() {
		pool.Close()
	}()

	storage := newPgFacade(pool)

	app := service.NewApp(storage)
	app.Run()
}

func newPgFacade(pool *pgxpool.Pool) *postgres.PgFacade {
	txManager := postgres.NewTxManager(pool)
	pgRepository := postgres.NewPgRepository(txManager)
	return postgres.NewPgFacade(txManager, pgRepository)
}

func newPgxPool(ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, "postgres://postgres:postgres@localhost:5432/orderdb?sslmode=disable")
	if err != nil {
		return nil, err
	}
	return pool, nil
}
