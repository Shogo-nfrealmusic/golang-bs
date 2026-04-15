package main

import "fmt"

func incrementGenerator() (func() int) {
	x := 0
	return func() int {
		x++
		return x
	}
}

func circleArea(pi float64) func(radius float64) float64 {
	return func(radius float64) float64 {
		return pi * radius * radius
	}
}

func main() {
	counter := incrementGenerator()
	fmt.Println(counter())
	fmt.Println(counter())
}