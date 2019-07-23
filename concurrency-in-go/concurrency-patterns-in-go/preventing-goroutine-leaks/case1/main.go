// ゴルーチンは生成コストは小さいが、ランタイムによってガベージコレクションされないため
// プロセス内に残留するとメモリリークにつながる可能性がある
package main

import (
	"fmt"
	"time"
)

func main(){
	Good()
}

func Baou() {
	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(completed)
			for s := range strings {
				// なにか面白い処理
				fmt.Println(s)
			}
		}()
		return completed
	}

	// nil をdoWorkに渡しているので、stringsチャネルには実際には文字列が書き込まれることはなく
	// doWorkを含むゴルーチンはこのプロセスが生きている限りずっとメモリ内に残る
	doWork(nil)
	fmt.Println("Done.")
}

// ゴルーチンの親子間で親から子にキャンセルのシグナルを遅れるようにする。
// 慣習としてこのシグナルは、doneという名前の読み込み専用チャネルにする
// 親ゴルーチンはこのチャネルを子ゴルーチンに渡して、キャンセルさせたいときにチャネルを閉じる
func Good() {
	// doneチャネルをdoWork関数に渡す。慣習としてこのチャネルは第1引数とする
	doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer  close(terminated)
			// for-selectパターン。case文の一つでdoneチャネルからシグナルが送られたか確認
			for {
				select {
				case s := <-strings:
					// なにか面白い処理
					fmt.Println(s)
				case <- done:
					return
				}
			}
		}()
		return terminated
	}
	done := make(chan interface{})
	terminated := doWork(done, nil)

	// i秒経過したらdoWorkの中で生成されたゴルーチンをキャンセルする他のゴルーチンを生成する
	go func() {
		// 1秒後に操作をキャンセルする
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	<-terminated
	fmt.Println("Done")
}
