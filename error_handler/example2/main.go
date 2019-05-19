package main

import "fmt"

const (
	NICK_MOCHO = "もちょ"
	NICK_TEN   = "天"
	NICK_NANSU = "ナンス"
)

// エラー構造体
type MochoError struct {
	Msg string
}

func (e *MochoError) Error() string {
	return "えらーだよー"
}

type NansuError struct {
	Msg  string
	Code int
}

func (e *NansuError) Error() string {
	return "なーんなーん"
}

func main() {
	printError(doError(NICK_MOCHO))
	printError(doError(NICK_TEN))
	printError(doError(NICK_NANSU))
}

func doError(nick string) error {
	if nick == NICK_MOCHO {
		return &MochoError{Msg: nick}
	}
	if nick == NICK_NANSU {
		return &NansuError{Msg: nick, Code: 19}
	}
	return nil
}

func printError(err error) {
	if err != nil {
		// 型でエラーの種類を判別する
		switch e := err.(type) {
		case *MochoError:
			fmt.Println("MochoError", err, "Msg:", e.Msg)
		case *NansuError:
			fmt.Println("NansuError:", err, "Msg:", e.Msg, "Code:", e.Code)
		default:
			fmt.Println("その他のエラー", err)
		}
		// if文の場合はこのように判定
		if nansu, ok := err.(*NansuError); ok {
			fmt.Println("NansuError[if", err, "Msg:", nansu.Msg, "Code:", nansu.Code)
		}
	} else {
		fmt.Println("no error")
	}
}
