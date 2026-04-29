package main

import (
	"fmt"
	"sync"
)

func player(name string, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		val, ok := <-ch
		if !ok {
			return
		}

		if val > 10 {
			close(ch)
			return
		}

		fmt.Println(name+":", val)
		ch <- val + 1
	}
}

func main() {
	ch := make(chan int)
	var wg sync.WaitGroup

	wg.Add(2)

	go player("Ping", ch, &wg)
	go player("Pong", ch, &wg)

	ch <- 1

	wg.Wait()
	fmt.Println("Game Over")
}
