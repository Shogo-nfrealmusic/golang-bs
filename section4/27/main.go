package main

import (
	"fmt"
	"time"
)

func getOsName() string {
	return "mac"
}
func main() {
	switch os := getOsName(); os {
	case "mac":
		fmt.Println("Mac OS")
	case "windows":
		fmt.Println("Windows")
	case "linux":
		fmt.Println("Linux")
	default:
		fmt.Println("Default")
	}

	t := time.Now()
	fmt.Println(t.Hour())
	switch {
	case t.Hour() < 12:
		fmt.Println("Morning")
	case t.Hour() < 17:
		fmt.Println("Afternoon")
	default:
		fmt.Println("Evening")
	}
}