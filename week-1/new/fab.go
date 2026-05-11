package main

import "fmt"

func main() {
	n := 20
	a := 0
	b := 1
	for i := 0; i < n; i++ {
		fmt.Print(a, " ")
		a, b = b, a+b
	}
}
