package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	ready := false
	for i := 1; i <= 10; i++ {
		go func(id int) {
			mu.Lock()
			for !ready {
				fmt.Println("Worker", id, "waiting...")
				cond.Wait()
			}
			fmt.Println("Worker", id, "started ")
			mu.Unlock()
		}(i)
	}
	time.Sleep(2 * time.Second)
	mu.Lock()
	fmt.Println("\nService is READY! Broadcasting...\n")
	ready = true
	cond.Broadcast()
	mu.Unlock()
	time.Sleep(2 * time.Second)
}
