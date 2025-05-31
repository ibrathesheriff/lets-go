package main

import (
	"fmt"
	"os"
)

func main() {
	var s string
	for _, arg := range os.Args {
		s += arg + "\n"
	}
	fmt.Println(s)
}
