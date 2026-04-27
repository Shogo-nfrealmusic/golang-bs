package main

import (
	"fmt"
	"golang-bs/section8/60/mylib"
	"golang-bs/section8/60/mylib/under"
)

func main() {
	s := []int{1, 2, 3, 4, 5}
	fmt.Println(mylib.Average(s))
	mylib.Say()
	under.Hello()
}