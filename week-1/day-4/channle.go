package main

import "fmt"

func main() {
	ch := make(chan int) // creating a channel of type int

	go func() {
		ch <- 100 // sending value to channel
	}()

	value := <-ch // receive value from channel
	fmt.Println(value)
}
