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

type Vertex3D struct {
	Vertex
	z int
}

func (v Vertex3D) Area3D() int {
	return v.x * v.y * v.z
}

func (v *Vertex3D) Scale3D(i int) {
	v.x = v.x * i
	v.y = v.y * i
	v.z = v.z * i
}

func New(x, y, z int) *Vertex3D {
	return &Vertex3D{Vertex{x, y}, z}
}

func main() {
	// v := Vertex{x: 10, y: 20}
	v := New(10, 20, 30)
	fmt.Println(v.Area3D())
}