package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Job struct {
	ID      int
	Retries int
}

func download(job Job) error {
	time.Sleep(100 * time.Millisecond)
	if rand.Intn(10) < 3 {
		return fmt.Errorf("failed")
	}
	return nil
}
func worker(id int, jobs chan Job, wg *sync.WaitGroup) {
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d (retry %d)\n", id, job.ID, job.Retries)
		err := download(job)
		if err != nil {
			if job.Retries < 3 {
				fmt.Printf("Retrying job %d\n", job.ID)
				job.Retries++
				go func(j Job) {
					jobs <- j
				}(job)
			} else {
				fmt.Printf(" Hard Failure:= job %d failed after 3 retries\n", job.ID)
				wg.Done()
			}
			continue
		}
		fmt.Printf(" Job %d completed\n", job.ID)
		wg.Done()
	}
}
func main() {
	rand.Seed(time.Now().UnixNano())
	numWorkers := 5
	numJobs := 20
	jobs := make(chan Job, 100)
	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		go worker(i, jobs, &wg)
	}
	for i := 1; i <= numJobs; i++ {
		wg.Add(1)
		jobs <- Job{ID: i, Retries: 0}
	}
	wg.Wait()
	close(jobs)
	fmt.Println(" All jobs processed including failures")
}
