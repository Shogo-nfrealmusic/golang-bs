package main

import "fmt"

func producer(first chan int) {
	for i := 0; i < 10; i++ {
		first <- i
	}
	close(first)
}

func multi2(first chan int, second chan int) {
	for i := range first {
		second <- i * 2
	}
	close(second)
}

func multi4(second chan int, third chan int) {
	for i := range second {
		third <- i * 4
	}
	close(third)
}

func main() {
	first := make(chan int)
	second := make(chan int)
	third := make(chan int)

	go producer(first)
	go multi2(first, second)
	go multi4(second, third)
	for result := range third {
		fmt.Println(result)
	}
}