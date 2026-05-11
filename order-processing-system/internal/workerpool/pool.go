package workerpool

import (
	"context"
	"errors"
	"log/slog"
	"sync"
)

type Job struct {
	ID      string
	Payload []byte
	Handle  func(context.Context) error
}

type Pool struct {
	jobs chan Job
	wg   sync.WaitGroup
	log  *slog.Logger
}

func New(size int, queue int, log *slog.Logger) *Pool {
	if size < 1 {
		size = 1
	}
	if queue < 1 {
		queue = size
	}
	return &Pool{jobs: make(chan Job, queue), log: log}
}

func (p *Pool) Start(ctx context.Context, workers int) {
	for i := 0; i < workers; i++ {
		p.wg.Add(1)
		go func(workerID int) {
			defer p.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-p.jobs:
					if !ok {
						return
					}
					if err := job.Handle(ctx); err != nil {
						p.log.Error("job failed", "worker_id", workerID, "job_id", job.ID, "error", err)
					} else {
						p.log.Info("job completed", "worker_id", workerID, "job_id", job.ID)
					}
				}
			}
		}(i + 1)
	}
}

func (p *Pool) Submit(ctx context.Context, job Job) error {
	if job.Handle == nil {
		return errors.New("job handler is required")
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case p.jobs <- job:
		return nil
	}
}

func (p *Pool) Shutdown() {
	close(p.jobs)
	p.wg.Wait()
}
