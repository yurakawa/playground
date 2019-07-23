package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	good()
}

// deferされたfmt.Printlnは実行されない。3回繰り返しが実行されたあとに
// 次の乱数の整数をもう読み込まれていないチャネルに送信しようとしてゴルーチンはブロックしてしまう
// case1のように生産者のゴルーチンに終了を伝えるチャネルを提供する
func bad() {
	newRandStream := func() <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			for {
				randStream <- rand.Int(
				)
			}
		}()
		return randStream
	}

	randStream := newRandStream()
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d,%d\n", i, <-randStream)
	}
}

func good() {
	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)

		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()
		return randStream
	}

	done := make(chan interface{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done)

	// 処理が実行中であることをシミュレート
	time.Sleep(1 * time.Second)
}
