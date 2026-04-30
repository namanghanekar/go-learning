package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, highPriority <-chan string, lowPriority <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-highPriority:
			if ok {
				fmt.Printf("Worker %d  HIGH: %s\n", id, job)
				time.Sleep(200 * time.Millisecond)
				continue
			}
			return
		default:
			select {
			case job, ok := <-highPriority:
				if ok {
					fmt.Printf("Worker %d  HIGH: %s\n", id, job)
				} else {
					return
				}
			case job, ok := <-lowPriority:
				if ok {
					fmt.Printf("Worker %d  LOW: %s\n", id, job)
				} else {
					return
				}
			}
			time.Sleep(200 * time.Millisecond)
		}
	}
}
func main() {
	highPriority := make(chan string, 5)
	lowPriority := make(chan string, 5)
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, highPriority, lowPriority, &wg)
	}
	for i := 1; i <= 5; i++ {
		lowPriority <- fmt.Sprintf("Low-%d", i)
	}
	for i := 1; i <= 5; i++ {
		highPriority <- fmt.Sprintf("High-%d", i)
	}
	close(highPriority)
	close(lowPriority)

	wg.Wait()
}
