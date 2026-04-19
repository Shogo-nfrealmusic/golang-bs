package main

import "fmt"

func main() {
	s := make([]int, 3)
	fmt.Printf("%T\n", s)

	m := make(map[string]int)
	fmt.Printf("%T\n", m)

	var p *int = new(int)
	fmt.Printf("%T\n", p)
}