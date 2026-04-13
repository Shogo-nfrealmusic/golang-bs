package main

import "fmt"

func main() {
	var i int = 1
	var f64 float64 = 1.2
	xf64 := 1.2
	var s string = "test"
	var t bool = true
	var f bool = false
	fmt.Println(i, f64, s, t, f)
	fmt.Printf("%T\n", xf64)

	xi := 1
	xs := "test"
	xt := true
	xf := false
	fmt.Println(xi, xf64, xs, xt, xf)
	fmt.Printf("%T\n", xi)
}