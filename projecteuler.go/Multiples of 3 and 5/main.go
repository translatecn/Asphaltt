package main

import (
	"fmt"
)

func main() {
	integers := make([]bool, 1000)
	for i := 1; i < 1000; i++ {
		if i%3 == 0 || i%5 == 0 {
			integers[i] = true
		}
	}

	sum := 0
	for i := 1; i < 1000; i++ {
		if integers[i] {
			sum += i
		}
	}
	fmt.Println(sum)
}
