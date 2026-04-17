package main

import (
	"fmt"
	"net/http"
	"sync"
)

func fetchURL(url string, wg *sync.WaitGroup) { // function to fetch a URL and print its status
	defer wg.Done() //DEFER is used to ensure that the Done() method is called when the function completes, regardless of how it exits (normal return or panic). This helps to prevent goroutine leaks and ensures that the WaitGroup counter is properly decremented.

	resp, err := http.Get(url) //http.Get is used to send an HTTP GET request to the specified URL and returns the response and any error encountered.
	if err != nil {
		fmt.Println("Error fetching:", url, err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Fetched:", url, "| Status:", resp.Status) //Prints the URL and its HTTP status code.
}

func main() {
	urls := []string{
		"https://github.com/namanghanekar",
		"https://www.linkedin.com/in/naman-ghanekar-78b6932a2/",
		"https://www.zenqua.com/",
		"https://docs.google.com/document/d/1oeuT5s8klkiS31p2KD0xiUYwhF9wOQT68Bd1K2Di5ac/edit?tab=t.86skwu6iqwn",
		"https://github.com/namanghanekar/go-learning/tree/main/week-1",
	}

	var wg sync.WaitGroup //sync.WaitGroup is a synchronization primitive that can be used to wait for a collection of goroutines to finish executing. It provides methods to add, done, and wait for goroutines.

	for _, url := range urls {
		wg.Add(1)
		go fetchURL(url, &wg) // concurrent execution
	}

	wg.Wait() // wait for all goroutines
	fmt.Println("All URLs fetched")
}
