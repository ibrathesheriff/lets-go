package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) > 1 {
		n, _ := strconv.Atoi(os.Args[1])
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				// place the star if the i and j are equal for going right
				// place for i matching to the end for going left
				if (i == j) || (i == n - j - 1) {
					fmt.Print("*")
				} else {
					fmt.Print(" ")
				}
			}
			fmt.Println()
		}
	}
}