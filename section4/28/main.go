package main

import (
	"fmt"
	"os"
)

func main() {
	// defer fmt.Println("World")
	// fmt.Println("Hello")

	file, _ := os.Open("main.go")
	defer file.Close()
	data := make([]byte, 100)
	file.Read(data)
	fmt.Println(string(data))
}