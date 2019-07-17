// 拘束は情報をたった一つの平行プロセスからのみ得られることを確実にしてくれる考え方
// => 並行プログラムは暗黙的に安全で同期が必要なくなる
// 拘束にはアドホックとレキシカルの２つ
package main

import (
	"bytes"
	"fmt"
	"sync"
)

func main() {
	lexicalStructs()
}

// アドホック拘束 -- 拘束を規約によって行う
// スライスであるdataがloopData関数でもhandleDataチャネルに対する繰り返しでも利用できることがわかる
// 以下のコードは規約によってloopData関数のみからアクセスしているので、間違ったコードが混入する可能性がある
func adhoc() {
	data := make([]int, 4)

	loopData := func(handleData chan<- int) {
		defer close(handleData)
		for i := range data {
			handleData <- data[i]
		}
	}

	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}
}



// レキシカル拘束
// 以下のコードのように拘束することでresultチャネルは直接利用できなくしている

func lexical() {
	chanOwner := func() <-chan int {
		//チャネルをchanOwner関数のレキシカルスコープ内で初期化する。
		// これによりresultsチャネルへの書き込みができるスコープを制限している
		// => 言い換えると、このチャネルへの書き込み権限を拘束して、他のゴルーチンの書き込みを防いでいる
		results := make(chan int, 5)
		go func() {
			defer close(results)
			for i := 0 ; i <= 5 ; i++ {
				results <- i
			}
		}()
		return results
	}

	// int のチャネルの読み込み専用のコピーを受け取る。読み込み権限のみが必要であることを宣言することで、
	// consumer関数内でのこのチャネルに対する操作を読み込み専用に拘束する
	consumer := func(results <- chan int) {
		for result := range results {
			fmt.Printf("Received: %d\n", result)
		}
	}
	// チャネルへの読み込み権限を受け取って消費者に渡す。消費者は読み込み以外は何もしない
	results := chanOwner()
	consumer(results)
}


// 平行安全でないデータ構造を使った例
// printDataはdataスライスの宣言のあとにないので、直接アクセスできず、
// 引数としてbyteのスライスを渡してもらう必要がある
// 起動したゴルーチンがそれぞれdataの一部しかアクセスできないように拘束している
// printDataを呼び出すゴルーチンでそれぞれ別の部分集合を渡す形式なのでメモリアクセスの同期や通信によるデータの共有が不要

// メリット: パフォーマンスの向上開発者に対する可読性の向上
//          => 同期のコストがなくなる
//
func lexicalStructs() {
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		var buff bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buff, "%c", b)
		}
		fmt.Println(buff.String())
	}

	var wg sync.WaitGroup
	wg.Add(2)
	data := []byte("golang")
	go printData(&wg, data[:3])
	go printData(&wg, data[3:])

	wg.Wait()
}
