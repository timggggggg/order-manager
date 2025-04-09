package outbox

import (
	"context"
	"log"
	"sync"
	"time"
)

type OutboxWorkerPool struct {
	outbox       *Outbox
	interval     time.Duration
	workersCount int
	wg           sync.WaitGroup
}

func NewOutboxWorkerPool(workersCount int, outbox *Outbox, interval time.Duration) (*OutboxWorkerPool, error) {
	return &OutboxWorkerPool{
		outbox:       outbox,
		interval:     interval,
		workersCount: workersCount,
	}, nil
}

func (w *OutboxWorkerPool) Start(ctx context.Context) {
	for range w.workersCount {
		w.wg.Add(1)
		go w.run(ctx)
	}
}

func (w *OutboxWorkerPool) Shutdown() {
	w.outbox.CloseProducer()
	w.wg.Wait()
}

func (w *OutboxWorkerPool) run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)

	defer w.wg.Done()
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.outbox.ProcessBatch(ctx)
		case <-ctx.Done():
			log.Println("Stopping outbox worker")
			return
		}
	}
}
