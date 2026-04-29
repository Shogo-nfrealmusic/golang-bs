package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	content, err := os.ReadFile("main.go")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(content))

	if err := os.WriteFile("main2.go", content, 0666); err != nil {
		log.Fatalln(err)
	}
}
