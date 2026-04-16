package main

import (
	"fmt"
	"net/http"
	"sync"
)

func fetchURL(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching:", url, err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Fetched:", url, "| Status:", resp.Status)
}

func main() {
	urls := []string{
		"https://github.com/namanghanekar",
		"https://www.linkedin.com/in/naman-ghanekar-78b6932a2/",
		"https://www.zenqua.com/",
		"https://docs.google.com/document/d/1oeuT5s8klkiS31p2KD0xiUYwhF9wOQT68Bd1K2Di5ac/edit?tab=t.86skwu6iqwn",
		"https://github.com/namanghanekar/go-learning/tree/main/week-1",
	}

	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go fetchURL(url, &wg) // concurrent execution
	}

	wg.Wait() // wait for all goroutines
	fmt.Println("All URLs fetched")
}
