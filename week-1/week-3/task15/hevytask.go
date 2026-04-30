package main

import (
	"context"
	"fmt"
	"time"
)

func heavyWork(ctx context.Context) {
	fmt.Println("Heavy work started")
	for {
		select {
		case <-ctx.Done():
			fmt.Println(" Stopping work")
			return
		default:
			fmt.Println("Working")
			time.Sleep(500 * time.Millisecond)
		}
	}
}
func main() {
	ctx,
		cancel := context.WithCancel(context.Background())
	go heavyWork(ctx)
	fmt.Println("Press ENTER to stop")
	fmt.Scanln()
	cancel()
	time.Sleep(1 * time.Second)
	fmt.Println("Program exited")
}
