package main

import "fmt"

func main() {
	l := []string{"python", "go", "java"}
	for i := 0; i < len(l); i++ {
		fmt.Println(i, l[i])
	}

	for i, v := range l {
		fmt.Println(i, v)
	}
}