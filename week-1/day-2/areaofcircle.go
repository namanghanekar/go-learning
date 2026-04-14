package main

import "fmt"

type Shape interface {
	Area() float64
}
type Circle struct {
	radius float64
}

func (c Circle) Area() float64 {
	return 3.14 * c.radius * c.radius
}

type Rectangle struct {
	length float64
	width  float64
}

func (r Rectangle) Area() float64 {
	return r.length * r.width
}

func main() {

	c := Circle{radius: 5}

	r := Rectangle{length: 4, width: 6}

	var s Shape

	s = c
	fmt.Println("Circle Area:", s.Area())

	s = r
	fmt.Println("Rectangle Area:", s.Area())
}
