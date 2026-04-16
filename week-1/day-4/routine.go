package main

import (
	"fmt"
	"time"
)

func new() {
	fmt.Println("heeeeeee")

}

func main() {
	go new()
	go func() {
		fmt.Println("Hello from goroutine!")
	}()
	time.Sleep(7000 * time.Second)
}
