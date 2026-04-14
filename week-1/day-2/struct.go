package main

import "fmt"

// 1️⃣ Define Struct
type Student struct {
	name  string
	age   int
	marks float64
}

// 2️⃣ Method (Read Operation)
func (s Student) display() {
	fmt.Println("Name:", s.name)
	fmt.Println("Age:", s.age)
	fmt.Println("Marks:", s.marks)
}

// 3️⃣ Method (Update using Pointer)
func (s *Student) updateMarks(newMarks float64) {
	s.marks = newMarks
}

// 4️⃣ Function (Create)
func createStudent(name string, age int, marks float64) Student {
	return Student{name: name, age: age, marks: marks}
}

func main() {

	// ✅ CREATE
	s1 := createStudent("Sreyansh", 21, 85.5)
	s2 := Student{"Rahul", 22, 90}

	// ✅ READ
	fmt.Println("---- Student 1 ----")
	s1.display()

	fmt.Println("---- Student 2 ----")
	s2.display()

	// ✅ UPDATE
	s1.updateMarks(95.0)

	fmt.Println("---- After Update studet 1  ----")
	s1.display()

	// ✅ POINTER USAGE
	p := &s2
	p.age = 25 // directly updating using pointer

	fmt.Println("---- After Pointer Update ----")
	s2.display()

}
