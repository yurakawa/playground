package main

import "fmt"

func main() {
	fmt.Println(bubbleSort([]int{3, 6, 1, 5, 7, 8}))
}

func bubbleSort(l []int) []int {
	for i := 0; i < len(l)-1; i++ {
		for j := i + 1; j < len(l); j++ {
			if l[i] > l[j] {
				l[i], l[j] = l[j], l[i]
			}
		}
	}
	return l
}
