package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	lifegame()
}


const (
	width = 80
	height = 15
)

type Universe [][]bool
func NewUniverse() Universe {
	u := make(Universe, height)
	for i := range u {
		u[i] = make([]bool, width)
	}
	return u
}
func (u Universe) string() string {
	var b byte
	buf := make([]byte, 0, (width+1)*height )

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			b = ' '
			if u[y][x] {
				b = '*'
			}
			buf = append(buf, b)
		}
		buf = append(buf, '\n')
	}

	return string(buf)
}

func (u Universe) show () {
	fmt.Print("\033[H\033[2J", u.string())
	// fmt.Print("\x0c", u.string())
}

func (u Universe) seed() {
	for i := 0; i < (width * height / 4); i++ {
		u.set(rand.Intn(width), rand.Intn(height), true)
	}
}

func (u Universe) set(x, y int, b bool) {
	u[y][x] = b
}


// セルは生きている？
func (u Universe)alive(x, y int) bool {
	x = (x + width) % width
	y = (y + height) % height
	return u[y][x]
}

// 隣接セルを数える
func (u Universe)neighbors(x, y int) int {
	n := 0
	for v := -1; v <= 1; v++ {
		for h := -1; h <= 1; h++ {
			if !(v == 0 && h == 0) && u.alive(x+h, y+v) {
				n++
			}
		}
	}
	return n
}

func (u Universe) Next(x, y int) bool {
	n := u.neighbors(x, y)
	return n == 3 || n == 2 && u.alive(x, y)
}

func step(a, b Universe) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			b.set(x, y, a.Next(x, y))
		}
	}
}

func lifegame () {
	a,b  := NewUniverse(), NewUniverse()
	a.seed()
	for i := 0; i < 300; i++ {
		step(a, b)
		a.show()
		time.Sleep(time.Second / 30)
		a, b = b, a // Swap universes
	}

}
