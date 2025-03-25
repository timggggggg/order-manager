package logpipeline

import (
	"context"
	"io"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	"gitlab.ozon.dev/timofey15g/homework/logger"
)

type LogPipeline struct {
	dbWP      *WorkerPool
	stdoutWP  *WorkerPool
	inputChan chan Log
}

func NewLogPipeline(ctx context.Context, stdout io.Writer, pool *pgxpool.Pool) *LogPipeline {
	inputDBChan := make(chan Log, 5)
	stdinChan := make(chan Log, 5)

	dbPool := NewWorkerPool(2, 5, 500*time.Millisecond, logger.NewConsoleLogger(stdout))
	stdoutPool := NewWorkerPool(2, 5, 500*time.Millisecond, logger.NewDBLogger(pool))

	dbPool.Start(ctx, inputDBChan, stdinChan)
	stdoutPool.Start(ctx, stdinChan, nil)

	return &LogPipeline{
		dbWP:      dbPool,
		stdoutWP:  stdoutPool,
		inputChan: inputDBChan,
	}
}

func (p *LogPipeline) LogStatusChange(TS time.Time, ID int64, statusFrom, statusTo models.OrderStatus) {
	log := Log{
		TS,
		ID,
		statusFrom,
		statusTo,
	}
	p.inputChan <- log
}

func (p *LogPipeline) Shutdown() {
	p.dbWP.Shutdown()
	p.stdoutWP.Shutdown()
}
