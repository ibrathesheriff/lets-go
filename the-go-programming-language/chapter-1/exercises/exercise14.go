package main

import (
	"fmt"
	"bufio"
	"os"
)

// Modify dup2 to print the names of all files in which each duplicated line occurs

func main() {
	counts := make(map[string]int)
	infile := make(map[string]string)
	// remove if related to the standard input stream
	//if len(os.Args[1:]) == 0 {
	//	countLines(os.Stdin, counts)
	//} 
	for _, arg := range os.Args[1:] {
		f, err := os.Open(arg)
		if err != nil {
			fmt.Fprint(os.Stderr, "dup exercise: %v\n", err)
			continue
		}
		countLines(f, counts, infile, arg)
	}
	
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s in files: %s\n", n, line, infile[line])
		}
	}
}

func countLines(f *os.File, counts map[string]int, infile map[string]string, filename string) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		counts[input.Text()]++
		infile[input.Text()] += " " + filename
	}
	// Note: ignoring potential errors from input.Err()
}