package main

import (
	"encoding/json"
	"fmt"
)

type Peee struct {
	ID  int
	Org string
}

func main() {
	var p Peee
	fmt.Println("enter id")
	fmt.Scan(&p.ID)
	fmt.Println("enter org")
	fmt.Scan(&p.Org)
	jsonData, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Struct: %+v\n", p)
	fmt.Println("JSON:", string(jsonData))
}
