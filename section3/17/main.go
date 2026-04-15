package main

import "fmt"

func main() {
	m := map[string]int{"apple": 100, "banana": 200}
	fmt.Println(m)
	v, ok := m["apple"]
	fmt.Println(v, ok)
	v, ok = m["nothing"]
	fmt.Println(v, ok)

	m2 := make(map[string]int)
	m2["pc"] = 5000
	fmt.Println(m2)
}
