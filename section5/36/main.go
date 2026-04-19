package main

import "fmt"

type Vertex struct {
	x int
	y int
	s string
}

func main() {
	v := Vertex{x: 1, y: 2}
	fmt.Println(v)

	v2 := Vertex{x: 1}
	fmt.Println(v2)

	v3 := Vertex{1, 2, "test"}
	fmt.Println(v3)
}