package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

func main() {
	tr(os.Stdin, os.Stdout, os.Stderr)
}

func tr(src io.Reader, dst io.Writer, errDst io.Writer) error {
	cmd := exec.Command("tr", "a-z", "A-Z")

	// 実行するコマンド tr a-z A-Z
	stdin, _  := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	err := cmd.Start() // コマンドの実行を開始する
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		// コマンドの標準入力にsrcからコピーする
		_, err := io.Copy(stdin, src)
		if e, ok := err.(*os.PathError); ok && e.Err == syscall.EPIPE{
			// ignore EPIPE
		} else if err != nil {
			log.Println("failed to write to STDIN", err)
		}
		stdin.Close()
		wg.Done()
	}()
	go func() {
		// コマンドの標準出力をdstにコピーする
		io.Copy(dst, stdout)
		stdout.Close()
		wg.Done()
	}()
	go func() {
		// コマンドの標準エラー出力をerrDstにコピーする
		stderr.Close()
		wg.Done()
	}()
	// 標準入出力のI/Oを行うgoroutineがすべて終わるまで待つ
	wg.Wait()
	// コマンドの終了を待つ
	return cmd.Wait()
}
