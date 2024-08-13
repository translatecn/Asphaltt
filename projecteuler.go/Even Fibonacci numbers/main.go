package main

import (
	"fmt"
)

func main() {
	sum, cnt := 2, 0
	x, y := 1, 2
	for y < 4_000_000 {
		x, y = y, x+y

		cnt++
		if cnt == 3 {
			// fmt.Println(y)
			sum, cnt = sum+y, 0
		}
	}

	fmt.Println(sum)
}
