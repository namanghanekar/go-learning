package main

import (
	"fmt"
	"time"
)

func slowAPI(ch chan string) {
	time.Sleep(700 * time.Millisecond)
	ch <- "API Response "
}
func callAPI() (string, error) {
	ch := make(chan string)
	go slowAPI(ch)
	select {
	case res := <-ch:
		return res, nil
	case <-time.After(500 * time.Millisecond):
		return "", fmt.Errorf("not responding")
	}
}
func main() {
	res, err := callAPI()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}
