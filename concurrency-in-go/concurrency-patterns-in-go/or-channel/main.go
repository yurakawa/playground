
package main

import (
	"fmt"
	"time"
)

// 再帰とゴルーチンを使って一つのdoneチャネルにまとめる
func main() {
	var or func(channels ...<-chan interface{}) <-chan interface{}
	// チャネルの可変長引数のスライスを受け取り、1つのチャネルを返します。
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		case 0: // スライスが空の場合はnilチャネルを返す
			return nil
		case 1: // スライスが1つしか要素を持っていない場合は、その要素を返すだけ
			return channels[0]
		}

		orDone := make(chan interface{})
		go func() { // 再帰が発生する部分。ゴルーチンを作ってブロックすることなく作ったチャネルにメッセージを受け取れるようにする。
			defer close(orDone)
			switch len(channels) {
			// 再帰のやり方のせいで、orへの各再帰呼び出しは少なくとも2つのチャネルを持っている。ゴルーチンの数を制限するために2つしかチャネルがなかった場合の特別な条件を設定する
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				// スライスの3番目以降のチャネルから再帰的にorチャネルを作成して、そこからselectを行う。
				// この再帰関係はスライスの残りの部分をorチャネルに分解して、最初のシグナルが帰ってくる木構造を形成する
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()
		return orDone
	}

	// 実際に使う
	// 1secondを渡した sig が合成されたチャネル全体を閉じる
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start)) // チャネルへの読み込みまでにかかった時間を表示する
}



