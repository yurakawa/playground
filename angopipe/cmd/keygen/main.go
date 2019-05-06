package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// 暗号鍵の作成 -- 32バイトの乱数を作り、それをbase64で文字列にする
func main() {
	key := make([]byte, 32)
	io.ReadFull(rand.Reader, key)
	readableKey := base64.StdEncoding.EncodeToString(key)
	fmt.Println(readableKey)
}
