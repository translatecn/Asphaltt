package main

import (
	"fmt"
)

func main() {
	// const oneMillion = 20
	const oneMillion = 1_000_000
	max, res, collats := 0, 0, make(map[int]int, oneMillion)
	for i := 1; i < oneMillion; i++ {
		coll := countCollatz(i, collats)
		collats[i] = coll
		if coll > max {
			max, res = coll, i
		}
	}

	fmt.Println(res, max)
}

func countCollatz(n int, collats map[int]int) int {
	res := 1
	for ; n != 1; res++ {
		if n%2 == 0 {
			n = n / 2
		} else {
			n = 3*n + 1
		}

		if c, ok := collats[n]; ok {
			return res + c
		}
	}
	return res
}
