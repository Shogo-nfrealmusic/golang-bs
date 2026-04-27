package main

import (
	"fmt"
	"golang-bs/section8/62/mylib"
	"golang-bs/section8/62/mylib/under"
)

func main() {
	numbers := []int{1, 2, 3, 4, 5}
	fmt.Println(mylib.Average(numbers))
	mylib.Say()
	under.Hello()
	p := mylib.Person{Name: "Mike", Age: 20}
	fmt.Println(p)
}
