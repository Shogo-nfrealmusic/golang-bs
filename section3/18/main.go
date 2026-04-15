package main

import "fmt"

func main() {
	b := []byte{
		109, 121, 32,  // m y (space)
		110, 97, 109, 101, 32, // name (space)
		105, 115, 32, // is (space)
		115, 104, 111, 103, 111, // shogo
	}

	fmt.Println(b)
	fmt.Println(string(b))
}