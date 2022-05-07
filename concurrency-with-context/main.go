package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	exampleWithDeadline()
	//err := doSomeThingConcurrency(3)
	//if err != nil {
	//	panic(err)

	//}
}

func doSomeThingConcurrency(workerNum int) error {
	// 必要なコンテキストを生成
	ctx := context.Background()
	cancelCtx, cancel := context.WithCancel(ctx)

	// 正常完了時にコンテキストのリソースを解放
	defer cancel()

	// 複数のゴルーチンからエラーメッセージを集約するためにチャネルを用意する
	errCh := make(chan error, workerNum)
	// workerNum分の並行処理を行う
	wg := sync.WaitGroup{}
	for i := 0; i < workerNum; i++ {
		i := i
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			// エラーが発生すれば、キャンセル処理尾をこない、エラーメッセージを送信する
			if err := doSomeThingWithContext(cancelCtx, num); err != nil {
				cancel()
				errCh <- err
			}
			return
		}(i)
	}

	// 並行処理の終了を待つ
	wg.Wait()

	// エラーチャネルに入ったメッセージを取り出す
	close(errCh)
	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	// エラーが発生していれば、最初のエラーを返す
	if len(errs) > 0 {
		return errs[0]
	}
	// 正常終了
	return nil
}

// コンテキストを利用した何らかの処理をする関数
func doSomeThingWithContext(ctx context.Context, num int) error {
	// 処理に入るまえに、コンテキストの死活を確認する
	select {
	case <-ctx.Done():
		return ctx.Err()
	// コンテキストがまだキャンセルされていなければ、そのまま処理に進む
	default:
	}
	fmt.Println(num)
	return nil
}

// 5秒後に
func exampleWithDeadline() {
	fmt.Println("start: ", time.Now())

	ctx := context.Background()
	d := time.Now().Add(5 * time.Second)
	parentCtx, cancel := context.WithDeadline(ctx, d)
	defer cancel()

	d2 := time.Now().Add(10 * time.Second)
	childCtx, cancel2 := context.WithDeadline(parentCtx, d2)
	defer cancel2()

	nd := d.AddDate(0, 0, 1)
	select {
	case <-time.After(time.Until(nd)):
		fmt.Println("ここは通らない")
	// case <-parentCtx.Done(): // (5)
	//	fmt.Println(time.Now())
	// 	fmt.Println("parentCtx: ", parentCtx.Err())
	case <-childCtx.Done(): // (5)
		fmt.Println(time.Now())
		fmt.Println("childCtx: ", childCtx.Err())
	}
}

/*
func exampleWithDeadline() {
	ctx := context.Background()
	// 指定時刻を生成
	d := time.Date(2022, 12, 18, 0, 0, 0, 0, time.UTC)
	// 指定時刻にキャンセルされるコンテキストを生成する
	timerCtx, cancel := context.WithDeadline(ctx, d)
	defer cancel()

	// 指定時刻の1日語の時刻を生成する
	nd := d.AddDate(0, 0, 1)
	// 自国ndになった時か、timeCtxがキャンセルされた時か、どちらか先の方が実行される
	select {
	case <-time.After(time.Until(nd)):
		fmt.Println("2022/12/19 00:00になりました。")
	case <-timerCtx.Done():
		fmt.Println(timerCtx.Err())
	}
}
func exampleWithTimeout() {
	ctx := context.Background()
	// 期間を決めてWithTimeoutでコンテキストを生成
	d := 15 * time.Second
	timerCtx, cancel := context.WithTimeout(ctx, d)
	// リソース解放を忘れない
	defer cancel()

	// 10秒経った時か、timerCtxがキャンセルされた時か、どちらか先の方が実行される
	select {
	case <-time.After(10 * time.Second):
		fmt.Println("10秒経ちました")
	case <-timerCtx.Done():
		fmt.Println(timerCtx.Err())
	}

	// さらに10秒が経った時か、timerCtxがキャンセルされた時か、どちらか先の方が実行される
	select {
	case <-time.After(10 * time.Second):
		fmt.Println("10秒経ちました")
	case <-timerCtx.Done():
		fmt.Println(timerCtx.Err())
	}
}
*/
