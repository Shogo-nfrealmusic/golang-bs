package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name string `json:"name"`
	Age int `json:"age"`
	Nicknames []string `json:"nicknames"`
}

func main() {
	b := []byte(`{"name": "Mike", "age": 20, "nicknames": ["a", "b", "c"]}`)
	var p Person
	if err := json.Unmarshal(b, &p); err != nil {
		fmt.Println(err)
	}
	fmt.Println(p)

	v, _ := json.Marshal(p)
	fmt.Println(string(v))
}