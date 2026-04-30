package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	URL        string
	StatusCode int
	Error      error
}

func fetchURL(client *http.Client, url string) Result {
	resp, err := client.Get(url)
	if err != nil {
		return Result{URL: url, StatusCode: 0, Error: err}
	}
	defer resp.Body.Close()
	return Result{
		URL:        url,
		StatusCode: resp.StatusCode,
		Error:      nil,
	}
}
func main() {
	urls := []string{
		"https://google.com",
		"https://github.com",
		"https://stackoverflow.com",
		"https://golang.org",
		"https://reddit.com",
		"https://news.ycombinator.com",
		"https://amazon.com",
		"https://facebook.com",
		"https://twitter.com",
		"https://linkedin.com",
		"https://microsoft.com",
		"https://apple.com",
		"https://netflix.com",
		"https://bing.com",
		"https://yahoo.com",
		"https://cnn.com",
		"https://bbc.com",
		"https://wikipedia.org",
		"https://openai.com",
		"https://youtube.com",
	}
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resultsChan := make(chan Result)
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			result := fetchURL(client, u)
			resultsChan <- result
		}(url)
	}
	go func() {
		wg.Wait()
		close(resultsChan)
	}()
	success := 0
	failure := 0
	for res := range resultsChan {
		if res.Error != nil {
			fmt.Printf(" %s → Error: %v\n", res.URL, res.Error)
			failure++
			continue
		}
		fmt.Printf(" %s → %d\n", res.URL, res.StatusCode)
		if res.StatusCode >= 200 && res.StatusCode < 300 {
			success++
		} else {
			failure++
		}
	}
	fmt.Println("\nSummary:=")
	fmt.Println("Total:", len(urls))
	fmt.Println("Success:", success)
	fmt.Println("Failure:", failure)
}
