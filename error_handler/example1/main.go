package main

import (
	"errors"
	"fmt"
)

const (
	NICK_MOCHO = "もちょ"
	NICK_TEN   = "天"
	NICK_NANSU = "ナンス"
)

var (
	ErrMocho = errors.New("えらーだよー")
	ErrNansu = errors.New("なーんなーん")
)

func main() {
	printError(doError(NICK_MOCHO))
	printError(doError(NICK_TEN))
	printError(doError(NICK_NANSU))

}

func doError(nick string) error {
	if nick == NICK_MOCHO {
		// ErrMochoインスタンスを返却
		return ErrMocho
	}

	if nick == NICK_NANSU {
		return ErrNansu
	}
	return nil
}
func printError(err error) {
	if err != nil {
		if err == ErrMocho {
			// ErrMocho用のエラー処理
			fmt.Println("ErrMocho", err)
		} else if err == ErrNansu {
			// ErrNansu用のエラー処理
			fmt.Println("ErrNansu", err)
		} else {
			fmt.Println("その他のエラー", err)
		}
	} else {
		fmt.Println("no error")
	}
}
