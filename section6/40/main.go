package main

import "fmt"

type Vertex struct {
	x, y int
}

func (v Vertex) Area() int {
	return v.x * v.y
}

func (v *Vertex) Scale(i int) {
	v.x = v.x * i
	v.y = v.y * i
}

// func Area(v Vertex) int {
// 	return v.x * v.y
// }

func New(x, y int) *Vertex {
	return &Vertex{x: x, y: y}
}

func main() {
	// v := Vertex{x: 10, y: 20}
	// fmt.Println(Area(v))
	v := New(10, 20)
	v.Scale(10)
	fmt.Println(v.Area())
}