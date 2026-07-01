package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"modulegue/core/queue"
)

type Processor func(ctx context.Context, payload json.RawMessage) (any, error)

type WorkerPool struct {
	queue      *queue.Queue
	processors map[string]Processor
	workers    int
	logger     *log.Logger
}

func NewWorkerPool(queue *queue.Queue, workers int, logger *log.Logger) *WorkerPool {
	if workers <= 0 {
		workers = 1
	}
	return &WorkerPool{
		queue:      queue,
		processors: make(map[string]Processor),
		workers:    workers,
		logger:     logger,
	}
}

func (p *WorkerPool) Register(topic string, processor Processor) {
	p.processors[topic] = processor
}

func (p *WorkerPool) Start(ctx context.Context) *sync.WaitGroup {
	var wg sync.WaitGroup
	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			p.runWorker(ctx, workerID)
		}(i + 1)
	}
	return &wg
}

func (p *WorkerPool) runWorker(ctx context.Context, workerID int) {
	idleDelay := 250 * time.Millisecond

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		job, err := p.queue.ClaimNext(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				time.Sleep(idleDelay)
				continue
			}
			if p.logger != nil {
				p.logger.Printf("worker-%d claim job error: %v", workerID, err)
			}
			time.Sleep(time.Second)
			continue
		}

		processor, ok := p.processors[job.Topic]
		if !ok {
			_ = p.queue.MarkRetry(ctx, job, errors.New("no processor registered for topic "+job.Topic))
			continue
		}

		result, err := processor(ctx, job.Payload)
		if err != nil {
			if p.logger != nil {
				p.logger.Printf("worker-%d process job %s topic=%s error: %v", workerID, job.ID, job.Topic, err)
			}
			_ = p.queue.MarkRetry(ctx, job, err)
			continue
		}

		if err := p.queue.MarkCompleted(ctx, job.ID, result); err != nil && p.logger != nil {
			p.logger.Printf("worker-%d mark complete job %s error: %v", workerID, job.ID, err)
		}
	}
}
