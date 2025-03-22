package logpipeline

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type ILogger interface {
	LogStatusChange(ctx context.Context, ts time.Time, id int64, statusFrom, statusTo models.OrderStatus)
}

type Log struct {
	TS         time.Time
	ID         int64
	StatusFrom models.OrderStatus
	StatusTo   models.OrderStatus
}

type WorkerPool struct {
	workersCount int
	batchSize    int
	timeout      time.Duration
	logger       ILogger
	wg           sync.WaitGroup
	input        chan Log
	next         chan Log
}

func NewWorkerPool(workersCount, batchSize int, timeout time.Duration, logger ILogger) *WorkerPool {
	return &WorkerPool{
		workersCount: workersCount,
		batchSize:    batchSize,
		timeout:      timeout,
		logger:       logger,
	}
}

func (wp *WorkerPool) Start(ctx context.Context, input chan Log, next chan Log) {
	wp.input = input
	wp.next = next
	for i := 0; i < wp.workersCount; i++ {
		wp.wg.Add(1)
		go wp.runWorker(ctx)
	}
}

func (wp *WorkerPool) runWorker(ctx context.Context) {
	var batch []Log
	timer := time.NewTimer(wp.timeout)

	defer wp.wg.Done()
	defer timer.Stop()

	for {
		select {
		case log, alive := <-wp.input:
			if !alive {
				wp.processBatch(ctx, batch)
				return
			}
			batch = append(batch, log)
			if len(batch) >= wp.batchSize {
				wp.processBatch(ctx, batch)
				batch = nil
				timer.Reset(wp.timeout)
			}

		case <-timer.C:
			wp.processBatch(ctx, batch)
			batch = nil
			timer.Reset(wp.timeout)

		case <-ctx.Done():
			wp.processBatch(ctx, batch)
			return
		}
	}
}

func (wp *WorkerPool) processBatch(ctx context.Context, batch []Log) {
	if len(batch) == 0 {
		return
	}
	for _, b := range batch {
		wp.logger.LogStatusChange(ctx, b.TS, b.ID, b.StatusFrom, b.StatusTo)
	}
	fmt.Println("batch processed")
	if wp.next != nil {
		for _, b := range batch {
			wp.next <- b
		}
	}
	fmt.Println("batch gone to the next worker")
}

func (wp *WorkerPool) Shutdown() {
	close(wp.input)
	wp.wg.Wait()
}
