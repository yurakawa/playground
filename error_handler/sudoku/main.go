package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	ErrBounds = errors.New("out of bounds") // 境界外
	ErrDigit = errors.New("invalid digit") // 不正な数
)
type SudokuError []error

// Errorは、1個以上のエラーをカンマで区切って返す
func (se SudokuError) Error () string {
	var s []string
	for _, err := range se {
		s = append(s, err.Error())
	}
	return strings.Join(s, ", ")
}

const rows, columns, digits = 9,9,9
type Grid [rows][columns]int8

func main() {
	var g Grid
	err := g.Set(10, 0, 15)
	if err != nil {
		if errs, ok := err.(SudokuError); ok {
			fmt.Printf("%d error(s) occurred:\n", len(errs))
			for _, e := range errs {
				fmt.Printf("- %v\n", e)
			}
		}
		os.Exit(1)
	}
	fmt.Printf("%+v", g)
}

func (g *Grid) Set (row, column int, digit int8) error {
	var errs SudokuError
	if !inBounds(row, column) {
		errs = append(errs, ErrBounds)
	}

	if !validDigit(digit) {
		errs = append(errs, ErrDigit)
	}
	if len(errs) > 0 {
		return errs
	}

	g[row][column] = digit
	return nil
}

func inBounds (row, column int) bool {
	if row < 0 || row >= rows {
		return false
	}
	if column < 0 || column >= columns {
		return false
	}
	return true
}

func validDigit (digit int8) bool {
	if digit < 0 || digit >= digits {
		return false
	}
	return true
}
