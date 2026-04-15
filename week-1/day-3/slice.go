package main

import "fmt"

func main() {
	sli := []int{1, 2, 3, 4, 5}
	fmt.Println(sli)
	fmt.Println(sli[3])        //accessing element according to their index
	sli = append(sli, 6, 7, 8) // adding new elements to the slice
	fmt.Println(sli)           // before updation
	sli[2] = 8
	fmt.Println(sli)        // after updation
	for i, v := range sli { //accessing all the elements using range loop
		fmt.Printf("Index: %d, Value: %d\n", i, v)
	}

}
