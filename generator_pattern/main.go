package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for i := range Count(ctx, 1, 99) {
		fmt.Println(i)
	}
}

func Count (ctx context.Context, start, end int) <- chan int{
	ch := make(chan int)
	go func(ch chan<-int) {
		defer close(ch)
	loop:
		for i := start; i <= end; i++ {

			select {
			case <-ctx.Done():
				break loop
			default:
			}

			// 重たい処理
			time.Sleep(500 * time.Millisecond)
			ch <- i
		}
	}(ch)
	return ch
}
