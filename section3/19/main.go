package main

import "fmt"

func add(x int, y int) int {
return (x + y)
}

func cal(price, item int) (result int) {
	result = price * item
	return
}

func main() {
	r := add(10, 20)
	fmt.Println(r)

	calPrice := cal(100, 5)
	fmt.Println(calPrice)

	f := func() {
		fmt.Println("inner func")
	}
	f()
}