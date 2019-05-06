package main

import (
	"angopipe"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	gcm, err := angopipe.Prepare()
	if err != nil {
		fmt.Fprintf(os.Stdout, "")
	}

	nonce := make([]byte, 12)
	n, err := io.ReadFull(os.Stdin, nonce)
	if n != 12 || err != nil {
		fmt.Fprintf(os.Stderr, "can't read nonce: %v\n", err)
		os.Exit(1)
	}

	encrypted, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "decryption error: %v\n", err)
		os.Exit(1)
	}

	result, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "decrption error: %v\n", err)
		os.Exit(1)
	}
	os.Stdout.Write(result)
}
