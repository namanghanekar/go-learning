package main

import (
	"fmt"
	"net/url"
)

func main() {
	fmt.Println("Urlparsing")
	urlString := "https://www.youtube.com/watch?v=vu6ZQ-t1sUk&list=PLzjZaW71kMwSEVpdbHPr0nPo5zdzbDulm&index=26"
	fmt.Printf("url type: %T\n", urlString)
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}
	fmt.Println("Parsed URL:-", parsedURL)
	fmt.Printf("url type:- %T\n", urlString)
	fmt.Printf("scheme:- %s\n", parsedURL.Scheme)
	fmt.Printf("host:- %s\n", parsedURL.Host)
	fmt.Printf("path:- %s\n", parsedURL.Path)

}
