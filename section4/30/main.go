package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("main.go")
	if err != nil {
		log.Fatalln("Error")
	}
	defer file.Close()
	data := make([]byte, 100)
	count, err := file.Read(data)
	if err != nil {
		log.Fatalln("Error")
	}
	fmt.Println(count, string(data))
	fmt.Println(string(data))

	if err = os.Chdir("test"); err != nil {
		log.Fatalln("Error")
	}
	fmt.Println("Success")
}