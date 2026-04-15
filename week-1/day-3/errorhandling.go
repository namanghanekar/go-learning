package main

import "fmt"

func test() {

	// defer + recover
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from:", r) //first defer delay the exicution because of error comes and recover will hwlp to recover from that error and print the message
		}
	}()

	fmt.Println("Step 1: Start function")

	// panic occurs
	panic("Something went wrong!") // if we are using withou defer or recover than it will  stop the exicution of program and print the message but with defer and recover it will print the message and continue the exicution of program

	fmt.Println("Step 2: End function") // will NOT run
}

func main() {
	test()
	fmt.Println("Program continues...")
}
