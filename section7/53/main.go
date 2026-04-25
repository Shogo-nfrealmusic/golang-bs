package main

import (
	"fmt"
	"sync"
)

func producer(ch chan int, i int) {
	ch <- i * 2
}

func consumer(ch chan int, wq *sync.WaitGroup) {
	for i := range ch {
		fmt.Println("received", i * 1000)
		wq.Done()
	}
}	

func main() {
	var wq sync.WaitGroup
	ch := make(chan int)

	for i := 0; i < 10; i++ {
		wq.Add(1)
		go producer(ch, i)
	}
	go consumer(ch, &wq)
	wq.Wait()
	close(ch)
}