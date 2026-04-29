package main

import "fmt"

const (
	c1 = iota
	c2 = iota
	c3 = iota
)

func main() {
	fmt.Println(c1, c2, c3)
}