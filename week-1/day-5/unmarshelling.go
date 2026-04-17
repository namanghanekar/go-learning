package main

import (
	"encoding/json"
	"fmt"
)

type Neee struct {
	ID  int    `json:"id"`
	Org string `json:"org"`
}

func main() {
	jsonData := `{"id":1,"zenqua":"ZenQua"}`
	fmt.Println(jsonData)
	var neee Neee
	json.Unmarshal([]byte(jsonData), &neee)
	fmt.Printf("%+v", neee)
}
