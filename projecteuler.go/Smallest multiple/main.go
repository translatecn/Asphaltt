package main

import (
	"fmt"
)

func main() {
	prod := 1

	// factors := []int{9, 8, 7, 5}
	factors := []int{19, 17, 16, 13, 11, 9, 7, 5}
	for _, f := range factors {
		prod *= f
	}

	fmt.Println(prod)
}
