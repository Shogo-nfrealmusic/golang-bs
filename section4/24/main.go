package main

import "fmt"

func by2(num int) string {
	if num % 2 == 0 {
		return "OK"
	} else {
		return "NG"
	}
}

func main() {
	result := by2(10)
	fmt.Println(result)
	
	num := 9
	if num % 2 == 0 {
		fmt.Println("even")
	} else if num % 3 == 0 {
		fmt.Println("3の倍数")
	} else {
		fmt.Println("odd")
	}

	x, y := 10, 2

	if x > y {
		fmt.Println("x is greater than y")
	} else if x < y {
		fmt.Println("x is less than y")
	} else {
		fmt.Println("x is equal to y")
	}
}
