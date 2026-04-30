package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	id int
}

func worker(id int, jobs <-chan Job, quit <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				fmt.Printf("Worker %d exiting (jobs closed)\n", id)
				return
			}
			fmt.Printf("Worker %d processing job %d\n", id, job.id)
			time.Sleep(100 * time.Millisecond)
		case <-quit:
			fmt.Printf("Worker %d stopped (scale down)\n", id)
			return
		}
	}
}

func main() {
	jobs := make(chan Job, 100)
	var wg sync.WaitGroup

	minWorkers := 2
	maxWorkers := 10
	workerCount := 0

	quitChans := []chan struct{}{}
	for i := 0; i < minWorkers; i++ {
		quit := make(chan struct{})
		quitChans = append(quitChans, quit)

		wg.Add(1)
		go worker(i, jobs, quit, &wg)
		workerCount++
	}
	stop := time.After(5 * time.Second)

	go func() {
		id := 1
		for {
			select {
			case jobs <- Job{id: id}:
				id++
				time.Sleep(50 * time.Millisecond)

			case <-stop:
				fmt.Println("\nStopping job producer...\n")
				close(jobs)
				return
			}
		}
	}()

	lastEmpty := time.Now()
	for {
		queueLen := len(jobs)
		if queueLen > 50 && workerCount < maxWorkers {
			quit := make(chan struct{})
			quitChans = append(quitChans, quit)

			wg.Add(1)
			go worker(workerCount, jobs, quit, &wg)
			workerCount++

			fmt.Println("Scaling UP → workers:", workerCount)
		}
		if queueLen == 0 {
			if time.Since(lastEmpty) > 5*time.Second && workerCount > minWorkers {
				close(quitChans[len(quitChans)-1])
				quitChans = quitChans[:len(quitChans)-1]
				workerCount--

				fmt.Println("Scaling DOWN → workers:", workerCount)
			}
		} else {
			lastEmpty = time.Now()
		}
		if queueLen == 0 && len(jobs) == 0 {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	wg.Wait()
	fmt.Println("\nAll workers finished. Program exit ✅")
}
