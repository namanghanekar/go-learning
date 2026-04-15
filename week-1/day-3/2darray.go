package main

import "fmt"

func main() {
	arr := [3][3]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	fmt.Println(arr)
	arr[1][2] = 10 // updating element at row 1 and column 2
	fmt.Println(arr)

}
