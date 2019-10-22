package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

func main() {
	if err := _main(); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}

func _main() error {
	if len(os.Args) < 2 {
		return errors.New("prog [file1 file2 ...]")
	}

	// シグナル処理のおまじない
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// os.Args の最初の引数はこのコマンド名
	// それを除外した分が対象ファイル名
	channels, err := makeChannelsForFiles(os.Args[1:])
	if err != nil {
		return err
	}

	// 上記ループから得たチャンネルから動的にselect caseを作成する
	cases, err := makeSelectCases(channels...)
	if err != nil {
		return err
	}

	// selectを動的人作成・実行する
	go doSelect(cases)

	// シグナルを受け取るまでブロックし続ける
	// Ctrl-Cとうつと終了
	// TODO:ここまで作成したgoroutineを正しく抜けてリソース開放を行う
	select {
	case <-sigch:
		return nil
	}
}

func readFromFile(ch chan []byte, f *os.File) {
	defer close(ch) // すべて終わったらチャンネルを閉じる
	defer f.Close() // すべて終わったらファイルを閉じる

	buf := make([]byte, 4096)
	for {
		// 読み込めるデータが有ればそれをチャネルに渡す
		// ここではエラーが有ってもファイルを置い続ける
		// （そうしないとio.EOFを受け取ったら、それ以上tailができなくなってしまう）
		if n, err := f.Read(buf); err == nil {
			ch <- buf[:n]
		}
	}
}

func makeChannelsForFiles(files []string)([]reflect.Value, error) {
	cs := make([]reflect.Value, len(files))

	for i, fn := range files {
		// データを流す用のチャンネルを作り
		ch := make(chan []byte)

		// ファイルをオープン
		f, err := os.Open(fn)
		if err != nil {
			return nil, err
		}
		go readFromFile(ch, f)

		cs[i] = reflect.ValueOf(ch)
	}
	return cs, nil
}

// チャンネルが格納されたreflect.Valueの配列を使い
// 対応するreflect.SelectCaseを作成
func makeSelectCases(cs ...reflect.Value) ([]reflect.SelectCase, error) {
	// 与えられた分のchanの数だけreflect.SelectCaseを作成
	cases := make([]reflect.SelectCase, len(cs))
	for i, ch := range cs {
		// reflect.Valueの値がチャンネルがない場合はエラー
		if ch.Kind() != reflect.Chan {
			return nil, errors.New("argument must be a channel")
		}

		// チャンネルの場合はSelectCaseを作成
		cases[i] = reflect.SelectCase{
			Chan: ch,
			Dir: reflect.SelectRecv,
		}
	}
	return cases, nil
}

// いずれかのselect caseからデータが帰ってきたら、
// それを標準出力に出力するループを繰り返し実行する
func doSelect(cases []reflect.SelectCase) {
	for {
		if chosen, recv, ok := reflect.Select(cases); ok {
			fmt.Printf("\n=== %s ===\n%s", os.Args[chosen+1], recv.Interface())
		}
	}
}
