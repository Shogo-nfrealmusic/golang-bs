package main

import (
	"fmt"
	"sort"
)

func main() {
	i := []int{9, 0, 8, 4, -1}
	s := []string{"a", "b", "c"}
	p := []struct {
		Name string
		Age int
	}{
		{Name: "Mike", Age: 20},
		{Name: "Nancy", Age: 24},
		{Name: "Messi", Age: 30},
	}
	fmt.Println(i, s, p)
	sort.Ints(i)
	fmt.Println(i)
}