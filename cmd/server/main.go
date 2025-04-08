package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	"gitlab.ozon.dev/timofey15g/homework/internal/kafka"
	logpipeline "gitlab.ozon.dev/timofey15g/homework/internal/log_pipeline"
	"gitlab.ozon.dev/timofey15g/homework/internal/logger"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	"gitlab.ozon.dev/timofey15g/homework/internal/outbox"
	"gitlab.ozon.dev/timofey15g/homework/internal/service"
	"gitlab.ozon.dev/timofey15g/homework/internal/storage/postgres"
	storagecache "gitlab.ozon.dev/timofey15g/homework/internal/storage_cache"
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

	brokers := []string{"localhost:9092"}
	tasksTable := "tasks"
	topic := "logs"

	producer, err := kafka.NewKafkaProducer(brokers)
	if err != nil {
		panic(fmt.Errorf("error creating kafka producer: %w", err))
	}

	ob := outbox.NewOutbox(pool, tasksTable, producer, topic)

	outboxWorkerPool, err := outbox.NewOutboxWorkerPool(2, ob, 500*time.Millisecond)
	if err != nil {
		log.Fatalf("error creating outboxWorkerPool: %v", err)
	}

	consumerWorkerPool, err := kafka.NewConsumerWorkerPool(1, brokers, topic)
	if err != nil {
		log.Fatalf("error creating consumerWorkerPool: %v", err)
	}

	outboxWorkerPool.Start(ctx)
	consumerWorkerPool.Start(ctx)

	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan
		cancel()
		logPipeline.Shutdown()
		outboxWorkerPool.Shutdown()
		consumerWorkerPool.Shutdown()
		pool.Close()
		os.Exit(0)
	}()

	storage := newPgFacade(pool)

	app := service.NewApp(storage, ob)
	app.Run()
}

func newPgFacade(pool *pgxpool.Pool) *postgres.PgFacade {
	txManager := postgres.NewTxManager(pool)
	pgRepository := postgres.NewPgRepository(txManager)

	cacheSize, err := strconv.ParseInt(os.Getenv("CACHE_SIZE"), 10, 64)
	if err != nil {
		panic(errors.New("invalid env variable CACHE_SIZE"))
	}

	cacheType := os.Getenv("CACHE_TYPE")

	cacheStrat, err := storagecache.NewCacheStrategy(cacheType, cacheSize)
	if err != nil {
		panic(errors.New("invalid env variable CACHE_TYPE"))
	}

	cache := storagecache.NewCache(cacheStrat)

	return postgres.NewPgFacade(txManager, pgRepository, cache, time.Now)
}

func newPgxPool(ctx context.Context, connectionString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
