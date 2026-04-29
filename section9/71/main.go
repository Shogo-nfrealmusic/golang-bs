package main

import (
	"context"
	"fmt"
	"time"
)

func longProcess(ctx context.Context, ch chan string) {
	fmt.Println("run")
	time.Sleep(4 * time.Second)
	fmt.Println("finished")
	ch <- "result"
}

func main() {
	ch := make(chan string)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	go longProcess(ctx, ch)
	
	for {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err())
			return
		case <-ch:
			fmt.Println("success")
			return
	}
}
}