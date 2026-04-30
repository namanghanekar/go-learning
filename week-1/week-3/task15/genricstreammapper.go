package main

import (
	"fmt"
)

func Map[T any, R any](in <-chan T, fn func(T) R) <-chan R {
	out := make(chan R)
	go func() {
		defer close(out)
		for val := range in {
			out <- fn(val)
		}
	}()
	return out
}
func main() {
	in := make(chan int)
	out := Map(in, func(x int) int {
		return x * x
	})
	go func() {
		for i := 1; i <= 5; i++ {
			in <- i
		}
		close(in)
	}()
	for result := range out {
		fmt.Println(result)
	}
}
