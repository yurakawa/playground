package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	const defaultPort = 3000
	// envPort, _ := strconv.Atoi(os.Getenv("PORT"))
	var port int
	flags := flag.NewFlagSet("example", flag.ContinueOnError)
	flags.SetOutput(os.Stdout)
	flags.Usage = func() { fmt.Println("Usage: example PATH")}
	flags.IntVar(&port, "port", defaultPort, "port to use")
	flags.IntVar(&port, "p", defaultPort, "port to use(short)")
	// flag.IntVar(&port, "port", envPort, "port to use")

	var spaces strSliceValue
	flags.Var(&spaces, "spaces", "")

	if err := flags.Parse(os.Args[1:]); err != nil {
		// エラー処理
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("%+v\n", spaces)
	fmt.Printf("%d\n", port)
}


type strSliceValue []string

func (v *strSliceValue) Set (s string) error {
	strs := strings.Split(s, ",")
	*v = append(*v, strs...)
	return nil
}

func (v *strSliceValue) String() string {
	return strings.Join(*v, ",")
}
