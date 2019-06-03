package main

import (
	"fmt"
	"golang.org/x/tour/pic"
	"golang.org/x/tour/wc"
	"strings"
)

func Pic(dx, dy int) [][]uint8 {

	sy := make([][]uint8, dy)
	for y := range sy {
		sx := make([]uint8, dx)
		for x := range sx {
			//value := (x + y) / 2
			value := x * y
			//value := x ^ y
			sx[x] = uint8(value)
		}
		sy[y] = sx
	}

	return sy
}

func exercicesSlices() {
	pic.Show(Pic)
}

func WordCount(s string) map[string]int {
	words := strings.Split(s, " ")
	m := make(map[string]int)
	for _, word := range words {
		m[word] = m[word] + 1
	}
	return m
}

func exerciseWordCount() {
	wc.Test(WordCount)
}

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	current, next := 1, 0
	return func() int {
		current, next = next, current+next
		return current
	}
}

func exerciseFib() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
