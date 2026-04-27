package main

import (
	"fmt"
	"golang-bs/section8/61/mylib"
	"golang-bs/section8/61/mylib/under"
)

func main() {
	s := []int{1, 2, 3, 4, 5}
	fmt.Println(mylib.Average(s))
	mylib.Say()
	under.Hello()
	p := mylib.Person{Name: "Mike", Age: 20}
	fmt.Println(p)
}