package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	"gitlab.ozon.dev/timofey15g/homework/internal/service"
	"gitlab.ozon.dev/timofey15g/homework/internal/storage/postgres"
	logpipeline "gitlab.ozon.dev/timofey15g/homework/log_pipeline"
	"gitlab.ozon.dev/timofey15g/homework/logger"
)

func main() {
	err := godotenv.Load("cmd/server/.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	ctx, cancel := context.WithCancel(context.Background())

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	db := os.Getenv("POSTGRES_DB")
	sslMode := os.Getenv("SSL_MODE")

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, db, sslMode)

	pool, err := newPgxPool(ctx, connectionString)
	if err != nil {
		log.Fatal("error newPgxPool", err)
	}

	storage := newPgFacade(pool)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./cmd/server")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	keywords := viper.GetStringSlice("keywords")

	filterWriter := &models.RequiredWordsWriter{
		Writer:        os.Stdout,
		RequiredWords: keywords,
	}

	inputDBChan := make(chan logpipeline.Log, 5)
	stdinChan := make(chan logpipeline.Log, 5)

	dbPool := logpipeline.NewWorkerPool(2, 5, 500*time.Millisecond, logger.NewConsoleLogger(filterWriter))
	stdoutPool := logpipeline.NewWorkerPool(2, 5, 500*time.Millisecond, logger.NewDBLogger(pool))

	dbPool.Start(ctx, inputDBChan, stdinChan)
	stdoutPool.Start(ctx, stdinChan, nil)

	logPipeline := logpipeline.GetLogPipelineInstance()
	logPipeline.SetWorkerPools(dbPool, stdoutPool)
	logPipeline.SetInputChan(inputDBChan)

	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan
		cancel()
		logPipeline.Shutdown()
		pool.Close()
		os.Exit(0)
	}()

	app := service.NewApp(storage)
	app.Run()
}

func newPgFacade(pool *pgxpool.Pool) *postgres.PgFacade {
	txManager := postgres.NewTxManager(pool)
	pgRepository := postgres.NewPgRepository(txManager)
	return postgres.NewPgFacade(txManager, pgRepository)
}

func newPgxPool(ctx context.Context, connectionString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
