package main

import "fmt"

func main() {
	x := 10

	p := &x // pointer storing address of x

	fmt.Println("Value of x:", x)
	fmt.Println("Address of x:", &x)
	fmt.Println("Pointer p:", p)
}
