//go:build !solution

package main

import (
	"fmt"
	"os"
	s "strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	wordcounter := make(map[string]int)
	for _, filename := range os.Args[1:] {
		data, err := os.ReadFile(filename)
		check(err)
		words := s.Split(string(data), "\n")
		for _, word := range words{
			wordcounter[word]++
		}
	}
	for word, cnt := range wordcounter {
		if cnt >= 2 {
			fmt.Printf("%d\t%s\n", cnt, word)
		}
	}
}
