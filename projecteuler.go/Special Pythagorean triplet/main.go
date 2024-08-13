package main

import (
	"fmt"
)

func main() {
	var a, b, c int
	for a = 1; a < 1000; a++ {
		for b = 999; b > a; b-- {
			c = 1000 - a - b

			aa, bb, cc := a*a, b*b, c*c
			if aa+bb == cc || aa+cc == bb {
				fmt.Println(a * b * c)
				return
			}
		}
	}

	panic("Result not found")
}
