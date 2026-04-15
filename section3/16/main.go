package main

import "fmt"

func main() {
	n := make([]int, 3, 5)
	fmt.Printf("len=%d cap=%d value=%v\n", len(n), cap(n), n)
	n = append(n, 0, 0)
	fmt.Printf("len=%d cap=%d value=%v\n", len(n), cap(n), n)
	n = append(n, 0, 0, 0, 0, 0)
	fmt.Printf("len=%d cap=%d value=%v\n", len(n), cap(n), n)

	c := make([]int,0, 5)
	for i := 0; i < len(c); i++ {
		c = append(c, i)
		fmt.Println(c)
	}
	fmt.Println(c)
}