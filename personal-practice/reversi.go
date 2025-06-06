package main

import (
	"os"
	"fmt"
)

type Stack struct {
	items []string
}

func main() {
	var n int
	n = len(os.Args)
	if n > 1 {
		wordStack := Stack{}
		for _, word := range os.Args[1:] {
			wordStack.Push(word)
		}
		for !wordStack.IsEmpty() {
			fmt.Printf("%s ", wordStack.Pop())
		}
		fmt.Println()
	} else {
		fmt.Println("No arguments provided!")
	}
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack) Push(data string) {
	s.items = append(s.items, data)
}

func (s *Stack) Pop() string {
	if s.IsEmpty() {
		return ""
	}
	top := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return top
}