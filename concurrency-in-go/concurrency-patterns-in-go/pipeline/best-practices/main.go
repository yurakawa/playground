package main

import "fmt"

func main(){
	bestPractice()
}

// データの塊をいっぺんに処理するバッチ処理
// ここ →for _, v := range add(multiply(ints, 2),1) {
func Batch() {
	// 新しいスライスを作って各要素に掛けていく
	multiply := func(values []int, multiplier int) []int {
		multipleValues := make([]int, len(values))
		for i, v := range values {
			multipleValues[i] = v * multiplier
		}
		return multipleValues
	}

	// 新しいスライスを作って各要素に足していく
	add := func(values []int, additive int) []int {
		addedValues := make([]int, len(values))
		for i, v := range values {
			addedValues[i] = v + additive
		}
		return addedValues
	}

	// multiplyとaddをつなげる
	ints := []int{1,2,3,4}
	for _, v := range add(multiply(ints, 2),1) {
		fmt.Println(v)
	}

}

// ここの値を一つずつ処理するストリーム処理
// この辺 → fmt.Println(multiply(add(multiply(v, 2), 1), 22))
// メモリフットプリントはパイプラインの入力サイズまで小さくなるが、パイプラインをforループの本体に入れて
// パイプラインに値を送り込む重労働をrangeにさせている。
// => スケールの可能性が制限されたり、パイプラインをループごとにインスタンス化している
func Stream() {
	multiply := func(value, multiplier int) int {
		return value * multiplier
	}
	add := func(value, additive int) int {
		return value + additive
	}

	ints := []int{1,2,3,4}
	for _, v := range ints {
		fmt.Println(multiply(add(multiply(v, 2), 1), 22))
	}
}

// preventing-goroutine-leaksのパターンを利用
// => どの関数もチャネルを返していていくつかは追加のチャネルも受け取っている
func bestPractice() {
	// 個別の値の塊をチャネル上に流れるデータのストリームに変換する
	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		intStream := make(chan int, len(integers))
		go func (){
			defer close(intStream)
			for _, i := range integers {
				select {
				case <-done:
					return
				case intStream <- i:
				}
			}
		}()
		return intStream
	}
	multiply := func(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int{
		multipliedStream := make(chan int)
		go func() {
			defer close(multipliedStream)
			for i := range intStream {
				select {
				case <- done:
					return
				case multipliedStream <- i * multiplier:
				}
			}
		}()
		return multipliedStream
	}

	add := func(done <-chan interface{}, intStream<-chan int, additive int) <-chan int{
		addedStream := make(chan int)
		go func() {
			defer close(addedStream)
			for i := range intStream {
				select {
				case <- done:
					return
				case addedStream <- i + additive:
				}
			}
		}()
		return addedStream
	}

	done := make(chan interface{})
	defer close(done)
	intStream := generator(done, 1,2,3,4)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
}
