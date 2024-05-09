package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	//cancel()
	//deadline()
	// waitCancelDone()
	waitBySelect()

}

func cancel() {
	ctx, cancel := context.WithCancel(context.Background())
	child(ctx)
	cancel()
	child(ctx)
}

func child(ctx context.Context) {
	if err := ctx.Err(); err != nil {
		fmt.Println("キャンセルされた")
		return
	}
	fmt.Println("キャンセルされていない")
}

func deadline() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
	fmt.Println("deadline:", ctx.Err())
	time.Sleep(2 * time.Second)
	fmt.Println("deadline:", ctx.Err())
	defer cancel()
}

func waitCancelDone() {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	go func() {
		fmt.Println("別のゴルーチン")
	}()
	fmt.Println("STOP")
	<-ctx.Done()
	fmt.Println("そして時は動き出す。")
}

func waitBySelect() {
	ctx, cancel := context.WithCancel(context.Background())
	task := make(chan int)
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("DONE")
				return
			case i := <-task:
				fmt.Println("get", i)
			default:
				fmt.Println("キャンセルされていない")
			}
		}
	}()

	time.Sleep(time.Second)
	for i := 0; i < 5; i++ {
		task <- i
	}
	cancel()
}
