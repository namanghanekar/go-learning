package main

import "fmt"

func main() {
	arr := [5]int{1, 2, 3, 4, 5}
	fmt.Println(arr)
	fmt.Println(arr[3])     //accessing element according to their index
	fmt.Println(arr)        // before updation
	arr[2] = 8              // updating element according to their index
	fmt.Println(arr)        //after updation
	for i, v := range arr { //accessing all the elements using range loop
		fmt.Printf("Index: %d, Value: %d\n", i, v)
	}
	fmt.Println(len(arr)) // length of array
}
