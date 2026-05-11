package main

import "fmt"

func main() {
	var n int
	fmt.Print("Enter number: ")
	fmt.Scan(&n)
	if n < 2 {
		fmt.Println("Not Prime")
		return
	}
	isPrime := true
	for i := 2; i < n; i++ {
		if n%i == 0 {
			isPrime = false
			break
		}
	}
	if isPrime {
		fmt.Println("Prime")
	} else {
		fmt.Println("Not Prime")
	}
}
