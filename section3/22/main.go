package main

import "fmt"

func foo(params ...int16) {
	fmt.Println(len(params), params)
}

func main() {
	foo(10, 20)
	foo(10, 20, 30)
}