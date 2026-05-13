package worker

import (
	"log"
)

type Job struct {
	OrderID int
}

type WorkerPool struct {
	JobQueue chan Job
}

func NewWorkerPool(workerCount int) *WorkerPool {

	pool := &WorkerPool{
		JobQueue: make(chan Job, 100),
	}

	for i := 0; i < workerCount; i++ {
		go pool.startWorker(i)
	}

	return pool
}

func (p *WorkerPool) startWorker(id int) {

	for job := range p.JobQueue {

		log.Printf(
			"Worker %d processing order %d\n",
			id,
			job.OrderID,
		)
	}
}

func (p *WorkerPool) AddJob(job Job) {
	p.JobQueue <- job
}
