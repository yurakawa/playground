package main

import "fmt"

func main() {
	for i := range Count(1, 99) {
		fmt.Println(i)
	}
}

func Count (start, end int) <- chan int{
	ch := make(chan int)
	go func(ch chan int) {
		for i := start; i <= end; i++ {
			ch <- i
		}
	}(ch)

	return ch
}
