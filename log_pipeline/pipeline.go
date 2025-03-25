package logpipeline

import (
	"sync"
	"time"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

var once sync.Once
var instance *LogPipeline = nil

type LogPipeline struct {
	dbWP      *WorkerPool
	stdoutWP  *WorkerPool
	inputChan chan Log
}

func GetLogPipelineInstance() *LogPipeline {
	once.Do(func() {
		instance = NewLogPipeline()
	})
	return instance
}

func (p *LogPipeline) SetWorkerPools(dbWP, stdoutWP *WorkerPool) {
	p.dbWP = dbWP
	p.stdoutWP = stdoutWP
}

func (p *LogPipeline) SetInputChan(ch chan Log) {
	p.inputChan = ch
}

func NewLogPipeline() *LogPipeline {
	return &LogPipeline{
		dbWP:      nil,
		stdoutWP:  nil,
		inputChan: nil,
	}
}

func (p *LogPipeline) LogStatusChange(TS time.Time, ID int64, statusFrom, statusTo models.OrderStatus) {
	if p.inputChan == nil {
		return
	}
	log := Log{
		0,
		TS,
		ID,
		statusFrom,
		statusTo,
		"",
		"",
		"",
		0,
		"",
	}
	p.inputChan <- log
}

func (p *LogPipeline) LogRequest(ts time.Time, method, url, request_body string) {
	if p.inputChan == nil {
		return
	}
	log := Log{
		1,
		ts,
		0,
		models.StatusDefault,
		models.StatusDefault,
		method,
		url,
		request_body,
		0,
		"",
	}
	p.inputChan <- log
}

func (p *LogPipeline) LogResponse(ts time.Time, code int64, body string) {
	if p.inputChan == nil {
		return
	}
	log := Log{
		2,
		ts,
		0,
		models.StatusDefault,
		models.StatusDefault,
		"",
		"",
		"",
		code,
		body,
	}
	p.inputChan <- log
}

func (p *LogPipeline) Shutdown() {
	p.dbWP.Shutdown()
	p.stdoutWP.Shutdown()
}
