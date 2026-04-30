package main

import (
	"fmt"
	"sync"
)

func partialSum(arr []int, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	sum := 0
	for _, val := range arr {
		sum += val
	}
	ch <- sum
}
func main() {
	arr := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		arr[i] = i + 1
	}
	segmentSize := len(arr) / 4
	ch := make(chan int, 4)
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		start := i * segmentSize
		end := start + segmentSize
		wg.Add(1)
		go partialSum(arr[start:end], ch, &wg)
	}
	wg.Wait()
	close(ch)
	total := 0
	for val := range ch {
		total += val
	}
	fmt.Println("Total Sum:", total)
}
