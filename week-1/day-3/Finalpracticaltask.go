package main

import (
	"fmt"
)

// function to get marks
func getMarks(students map[string]int, name string) (int, error) {
	marks, ok := students[name]

	if !ok {
		return 0, fmt.Errorf("student not found")
	}

	return marks, nil
}

// function to display all students
func displayAll(students map[string]int) {
	fmt.Println("\nAll Students Data:")
	for name, marks := range students {
		fmt.Printf("Name: %s, Marks: %d\n", name, marks)
	}
}

func main() {
	// student data (map)
	students := map[string]int{
		"Naman":     90,
		"Priya":     85,
		"Shreyansh": 95,
	}

	var choice int

	fmt.Println("---- Student Grade System ----")
	fmt.Println("1. Search Student")
	fmt.Println("2. Fetch All Students")
	fmt.Print("Enter your choice: ")
	fmt.Scan(&choice)

	switch choice {

	case 1:
		var name string
		fmt.Print("Enter student name: ")
		fmt.Scan(&name)

		marks, err := getMarks(students, name) // calling getMarks function to get marks of student according to name

		if err != nil {
			fmt.Println("Error:", err) //if student not found in the map  than it will print the error message
			return
		}

		fmt.Printf("Marks of %s: %d\n", name, marks)

	case 2:
		displayAll(students)

	default:
		fmt.Println("Invalid choice")
	}
}
