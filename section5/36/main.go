package main

import "fmt"

type Vertex struct {
	x int
	y int
	s string
}

func changeVertex(v Vertex) {
	v.x = 1000
}

func changeVertex2(v *Vertex) {
	v.x = 1000
}

func main() {
	v := Vertex{x: 1, y: 2, s: "test"}
	changeVertex(v)
	fmt.Println(v)

	v2 := &Vertex{x: 1, y: 2, s: "test"}
	changeVertex2(v2)
	fmt.Println(*v2)
	/*
	v := Vertex{x: 1, y: 2}
	fmt.Println(v)

	v2 := Vertex{x: 1}
	fmt.Println(v2)

	v3 := Vertex{1, 2, "test"}
	fmt.Println(v3)

	var v5 Vertex
	fmt.Println(v5)

	v6 := new(Vertex)
	fmt.Printf("%T\n", v6)
	*/
	
}