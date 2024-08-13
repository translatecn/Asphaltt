package main

import (
	"fmt"
)

const (
	fiveHundred = 500
)

func main() {
	var divisors int
	triangle, i := 1, 2
	for ; divisors <= fiveHundred; triangle, i = triangle+i, i+1 {
		divisors = getDivisors(triangle)
	}
	fmt.Println(i-2, triangle-i+1)
}

func getDivisors(n int) (res int) {
	for i, j := 1, n; i <= j; i, j = i+1, n/(i+1) {
		if i*j != n {
			continue
		}

		if i != j {
			res += 2
		} else {
			res++
		}
	}
	return
}
