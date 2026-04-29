package main

import (
	"fmt"
	"strings"
	"unicode"
)

func generator(input []string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for _, str := range input {
			out <- str
		}
	}()
	return out
}
func sanitizer(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for str := range in {
			str = strings.ToLower(str)
			clean := strings.Map(func(r rune) rune {
				if unicode.IsPunct(r) {
					return -1
				}
				return r
			}, str)
			out <- clean
		}
	}()
	return out
}
func publisher(in <-chan string) {
	for str := range in {
		fmt.Println(str)
	}
}
func main() {
	input := []string{
		"HELLO",
		"My name is",
		"NAMAN",
	}
	gen := generator(input)
	san := sanitizer(gen)
	publisher(san)
}
