package main

import "fmt"

func main() {
	m := map[string]int{
		"naman":     1,
		"priya":     2,
		"shreyansh": 3, // we have a three key value pairs
	}
	value, ok := m["piyush"] //access a value according to key

	if ok {
		fmt.Println("Found:", value) //if present its return value
	} else {
		fmt.Println("Not found") //if not present in mao than return not found
	}
	for k, v := range m {
		fmt.Printf("Key: %s, Value: %d\n", k, v) //this is for access all the elements using range loop
	}
	fmt.Println(m)
	fmt.Println(m["naman"])

	m["naman"] = 10
	fmt.Println(m)
	delete(m, "priya") // this is for delete a key value pair from map
	fmt.Println(m)
}
