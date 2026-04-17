package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Item struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

var fileName = "inventory.json"

func loadData() map[int]Item {
	inventory := make(map[int]Item)

	file, err := os.ReadFile(fileName)
	if err == nil {
		json.Unmarshal(file, &inventory)
	}

	return inventory
}

func saveData(inventory map[int]Item) {
	data, _ := json.MarshalIndent(inventory, "", "  ")
	os.WriteFile(fileName, data, 0644)
}

func printJSON(data interface{}) {
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(jsonData))
}

func main() {
	inventory := loadData()

	for {
		fmt.Println("\n===== Inventory Menu =====")
		fmt.Println("1. Add Item")
		fmt.Println("2. Update Item")
		fmt.Println("3. Search Item by ID")
		fmt.Println("4. View All Items")
		fmt.Println("5. Remove Item")
		fmt.Println("6. Exit")

		var choice int
		fmt.Print("Enter your choice: ")
		fmt.Scan(&choice)

		switch choice {

		case 1:
			var id int
			var name string
			var price float64
			var quantity int

			fmt.Print("Enter ID: ")
			fmt.Scan(&id)
			if _, exists := inventory[id]; exists {
				fmt.Println("Item already exists!")
				break
			}

			fmt.Print("Enter Name: ")
			fmt.Scan(&name)

			fmt.Print("Enter Price: ")
			fmt.Scan(&price)

			fmt.Print("Enter Quantity: ")
			fmt.Scan(&quantity)

			inventory[id] = Item{name, price, quantity}
			saveData(inventory)

			fmt.Println("Item added!")

		case 2:
			var id int
			fmt.Print("Enter ID: ")
			fmt.Scan(&id)

			item, exists := inventory[id]
			if !exists {
				fmt.Println("Item not found!")
				break
			}

			fmt.Print("Enter New Name: ")
			fmt.Scan(&item.Name)

			fmt.Print("Enter New Price: ")
			fmt.Scan(&item.Price)

			fmt.Print("Enter New Quantity: ")
			fmt.Scan(&item.Quantity)

			inventory[id] = item
			saveData(inventory)

			fmt.Println("Item updated!")

		case 3:
			var id int
			fmt.Print("Enter ID: ")
			fmt.Scan(&id)

			item, exists := inventory[id]
			if exists {
				printJSON(item)
			} else {
				fmt.Println("Item not found!")
			}

		case 4:
			printJSON(inventory)

		case 5:
			var id int
			fmt.Print("Enter ID: ")
			fmt.Scan(&id)

			if _, exists := inventory[id]; exists {
				delete(inventory, id)
				saveData(inventory)
				fmt.Println("Item removed!")
			} else {
				fmt.Println("Item not found!")
			}

		case 6:
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid choice!")
		}
	}
}
