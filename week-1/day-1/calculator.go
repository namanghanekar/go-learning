package main

import (
	"fmt"
)

func main() {
	var a, b float64
	var op string

	fmt.Print("Enter calculation : ")

	fmt.Scan(&a, &op, &b) // taking input from user

	switch op {
	case "+":
		fmt.Println("Result:", a+b)
	case "-":
		fmt.Println("Result:", a-b)
	case "*":
		fmt.Println("Result:", a*b)
	case "/":
		if b == 0 {
			fmt.Println("Cannot divide by zero")
			return
		}
		fmt.Println("Result:", a/b)
	default:
		fmt.Println("Invalid operator")
	}
}
