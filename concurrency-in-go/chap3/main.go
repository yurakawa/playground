package chap3

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"runtime"
	"sync"
	"text/tabwriter"
	"time"
)

func main() {
	selectSample4()
}

func selectSample4() {
	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	workCounter := 0
loop:
	for {
		select {
		case <-done:
			break loop
		default:
		}

		// Simulate work
		workCounter++
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("Achieved %v cycles of work before signalled to stop .\n", workCounter)

}

func selectSample3 (){
	var c <- chan int
	start := time.Now()
	select {
	case <-c:
		fmt.Println("...")
	case <- time.After(1 * time.Second):
		fmt.Println("Timed out.")
	default:
		fmt.Printf("In default after %v\n\n", time.Since(start))
	}

}
func selectSample2() {
	c1 := make(chan interface{}); close(c1)
	c2 := make(chan interface{}); close(c2)

	var c1Count, c2Count int
	for i := 1000; i >= 0 ;i-- {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}
	fmt.Println("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)
}

func selectSample() {
	start := time.Now()
	c :=make(chan interface{})
	go func() {
		time.Sleep(5*time.Second)
		close(c) // 5秒待ったあとにチャネルを閉じる
	}()

	fmt.Println("Blocking on read...")
	select {
	case <-c: // チャネルの読み込みを試す
		fmt.Printf("Unblocked %v later.\n", time.Since(start))
	}
}

func chanCapsuelSample() {
	chanOwner := func() <-chan int {
		resultCh := make(chan int, 5) // バッファ付きチャネルを初期化。結果は６つ生成するとわかっているのでゴルーチンができる限り早く完了するようにキャパシティが5のバッファ付きチャネルを作成する
		go func() { // resultChへ書き込みを行うための無名ゴルーチンを起動。ゴルーチンより先にチャネルを生成したことに注意。外の関数によってカプセル化されている
			defer close(resultCh) //resultChを使い終わったあとに確実に閉じられるようにしておく。チャネルの所有者としての責任
			for i:= 0; i <=5; i++ {
				resultCh <- i
			}
		}()
		return resultCh // チャネルを返す。戻り値は読み込み専用として宣言されているので、resultChは暗黙的に読みお店尿の消費者に変換される
	}

	resultCh := chanOwner()
	for result := range resultCh { // resultChをrangeで繰り返す。消費者としてチャネルのブロックとチャネルを閉じたことだけに注意する
		fmt.Printf("Received: %d\n", result)
	}
	fmt.Println("Done receiving!")
}

func chanWithBufferSample() {
	var stdoutBuff bytes.Buffer // インメモリのバッファを作って出力が非決定になるのを軽減する。stdoutに直接書き込むより若干早い
	defer stdoutBuff.WriteTo(os.Stdout) // プロセスが終了する前に確実にバッファがstdoutに書きっこまれるようにする

	intCh := make(chan int, 4)
	go func() {
		defer close(intCh)
		defer fmt.Fprintln(&stdoutBuff, "Producer Done.")
		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
			intCh <- i
		}
	}()

	for integer := range intCh {
		fmt.Fprintf(&stdoutBuff, "Received %v.\n", integer)
	}
}

func chanReleaseSample() {
	begin := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin //ここでチャネル読み込み可能までゴルーチンが待機
			fmt.Printf("%v has begun\n", i)
		}(i)
	}

	fmt.Println("Unblocking goroutines...")
	close(begin) // これによってすべてのゴルーチンを同時に開放する
	wg.Wait()
}

func chanRoopSample() {
	intCh := make(chan int)
	go func() {
		defer close(intCh)
		for i:= 1; i<=5; i++ {
			intCh <- i
		}
	}()

	for integer := range intCh {
		fmt.Printf("%v ", integer)
	}
}

func chanSample2() {
	stringCh := make(chan string)
	go func() {
		stringCh <- "Hello channels!"
	}()
	close(stringCh)
	salutation, ok := <-stringCh
	fmt.Printf("(%v): %v", ok, salutation)
}

func chanSample() {
	stringCh := make(chan string)
	go func() {
		stringCh <- "Hello channels!"
	}()
	fmt.Println(<-stringCh)
}


func poolSample2() {
	var numCalcsCreated int
	calcPool := &sync.Pool {
		New: func() interface{} {
			numCalcsCreated += 1
			mem := make([]byte, 1024)
			return &mem // バイトのスライスのアドレスを保存していることに注意
		},
	}

	// プールに4kb確保する
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())

	const numWorkers = 1024*1024
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()

			mem := calcPool.Get().(*[]byte) // バイトのスライスのポインタであると型アサーションする
			defer calcPool.Put(mem)

			// 処理
		}()
	}
	wg.Wait()
	fmt.Printf("%d calculators were created.", numCalcsCreated)
}

func poolSample() {
	myPool := &sync.Pool {
		New: func() interface{} {
			fmt.Println("Creating new instance.")
			return struct{}{}
		},
	}
	myPool.Get() // プールのGetを呼び出す。下との2つの呼び出しによって、プールに定義されているNew関数を起動する。(インスタンスが初期化されていないから)
	instance := myPool.Get()
	myPool.Put(instance) // 先にプールから取得したインスタンスをプールに戻す。これで利用できるインスタンスの数を1に増やす
	myPool.Get() // この呼出が実行されたときは、先に生成されてプールに戻されたインスタンスを再利用する(Newが呼び出されない)
}

func onceSample() {
	var count int
	increment := func() {
		count++
	}

	var once sync.Once
	var increments sync.WaitGroup
	increments.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer increments.Done()
			once.Do(increment)
		}()
	}
	increments.Wait()
	fmt.Printf("Count is %d\n", count)
	// => 1
}

func onceSample2() {
	var count int
	increment := func() { count++ }
	decrement := func() { count-- }
	var once sync.Once
	once.Do(increment)
	once.Do(decrement)

	fmt.Printf("Count: %d\n", count)
	// => 1
}

func condBroadcastSample() {
	type Button struct { // Clicked という条件を含んでいるButton型を定義する
		Clicked *sync.Cond
	}

	button := Button { Clicked: sync.NewCond(&sync.Mutex{}) }

	 //条件に応じて送られてくるシグナルを扱う関数を登録するための便利な関数を定義する。
	 // 各ハンドラーはそれぞれのゴルーチン上で動作する。subscribeはゴルーチンが実行されていると確認できるまで終了しない
	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)
		go func() {
			goroutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		goroutineRunning.Wait()
	}

	var clickRegistered sync.WaitGroup // プログラムがstdoutへ書き込む前に終了してしまわないようにするためのWaitGroup
	clickRegistered.Add(3)
	subscribe(button.Clicked, func() { // ボタンがクリックされたときにボタンが有るウインドウを最大化するのをシミュレートしたハンドラを登録
		fmt.Println("Maximizing window.")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() { // マウスがクリックされたときにダイアログボックスを表示するのをシミュレートしたハンドラーを登録
		fmt.Println("Displaying annoying dialog box!")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() { // ユーザがアプリケーションのボタンをクリックした状態からマウスのボタンを離した状態をシミュレートする
		fmt.Println("Mouse clicked.")
		clickRegistered.Done()
	})
	// マウスのボタンが離されたときのハンドラーを設定する。こちらはClickedという状態(Cond) に対応するBroadcastを呼び出して、
	// すべてのハンドラーにマウスのボタンがクリックしたということを知らせる。
	button.Clicked.Broadcast()

	clickRegistered.Wait()
}

func condSample() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10) // 長さ0のスライスを作成。。最終的に10買い足すとわかっているのでキャパシティを10に設定

	removeFromQueue := func(dilay time.Duration) {
		time.Sleep(dilay)
		c.L.Lock() // 再度クリティカルセクションに入り条件にあった形でデータを修正する
		queue = queue[1:] // スライスと先頭をスライスの2番目の要素を指すように変えることで急ーから取り出したことにする
		fmt.Println("Removed from queue")
		c.L.Unlock() // キューから要素を取り出したので条件のクリティカルセクションを抜ける
		c.Signal() // 条件を待っているゴルーチンに何かが起きたことを知らせる
	}

	for i := 0; i < 10; i++ {
		c.L.Lock() // 条件であるLockerのLockメソッドを呼び出してクリティカルセクションに入る
		for len(queue) == 2 { // ループ内でキューの長さを確認する。条件上のシグナルは必ずしも同じ待っている事象が起きたことを意味しないので、この条件追加が必要
			c.Wait() // 条件のシグナルが創出されるまでメインゴルーチンを一時停止する
		}
		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second) // 1秒後に要素をキューから取り出す新しいゴルーチンを生成する
		c.L.Unlock()
	}
}

func condMinimum() {
	c := sync.NewCond(&sync.Mutex{}) // 新しいCondのインスタンスを作成
	c.L.Lock() // この条件でLockerをロックする。Waitへの呼び出しがループに入る時に自動的にUnlockを呼び出すため
	for conditionTrue() == false {
		c.Wait() // 条件が発生したかどうか通知を待つ。これはブロックする呼び出しでゴルーチンは一時停止する
	}
	c.L.Unlock() //この条件でLockerのロックを解除する。この記述はWaitの呼び出しが終わるとこの条件でLockを呼び出すので必要
}

func conditionTrue() bool {
	return true
}

// rwMutexSample ...
func rwMutexSample() {
	producer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		l.Lock()
		l.Unlock()
		time.Sleep(1)
	}

	observer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		l.Lock()
		defer l.Unlock()
	}

	test := func(count int, mutex, rwMutex sync.Locker) time.Duration {
		var wg sync.WaitGroup
		wg.Add(count+1)
		beginTestTime := time.Now()
		go producer(&wg, rwMutex)
		for i := count; i > 0; i-- {
			observer(&wg, rwMutex)
		}
		wg.Wait()
		return time.Since(beginTestTime)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer tw.Flush()

	var m sync.RWMutex
	fmt.Fprintf(tw, "Readers\tRWMutex\tMutex\no")
	for i := 0; i < 20; i++ {
		count := int(math.Pow(2, float64(i)))
		fmt.Fprintf(
			tw,
			"%d\t%v\t%v\n",
			count,
			test(count, &m, m.RLocker()), // 読み取りロックの読み取り
			test(count, &m, &m), // 書き込みロックの読み取り
		)
	}

}

func mutexSample() {
	var count int
	var lock sync.Mutex

	increment := func() {
		lock.Lock()
		defer lock.Unlock()
		count ++
		fmt.Printf("Incrementing: %d\n",count)
	}

	decrement := func() {
		lock.Lock()
		defer lock.Unlock()
		count--
		fmt.Printf("Decrementing: %d\n", count)
	}

	// インクリメント
	var arithmetic  sync.WaitGroup
	for i := 0; i <= 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			increment()
		}()
	}

	// デクリメント
	for i := 0; i <= 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			decrement()
		}()
	}
	arithmetic.Wait()
	fmt.Println("Arithmetic complete.")
}

func waitGroupSample() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("1st goroutine sleeping...")
		time.Sleep(1)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("2nd goroutine sleeping...")
		time.Sleep(2)
	}()

	wg.Wait()
	fmt.Println("All goroutines complete.")
}

func waitGroupForRoop() {
	hello := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		fmt.Printf("Hello from %v!\n", id)
	}

	const numGreeters = 5
	var wg sync.WaitGroup
	wg.Add(numGreeters)
	for i := 0; i < numGreeters; i++ {
		go hello(&wg, i + 1)
	}
	wg.Wait()
}

func memoryConsume () {
	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}
	var c <- chan interface{}
	var wg sync.WaitGroup
	noop := func() { wg.Done(); <-c }

	const numGoroutines = 1e4
	wg.Add(numGoroutines)
	before := memConsumed()
	for i := numGoroutines; i > 0; i-- {
		go noop()
	}
	wg.Wait()
	after := memConsumed()
	fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1000)
}



