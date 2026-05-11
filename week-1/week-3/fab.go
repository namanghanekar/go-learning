package main

import "fmt"

func main() {
	n := 10 // number of terms

	a, b := 0, 1

	for i := 0; i < n; i++ {
		fmt.Print(a, " ")
		a, b = b, a+b
	}
}
