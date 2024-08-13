package main

import (
	"fmt"
)

func main() {
	const topline = 1000
	num := make([]byte, topline)
	num[0] = 1

	for i := 1; i <= topline; i++ {
		extra := byte(0)
		for j := 0; j < topline; j++ {
			n := num[j]*2 + extra
			num[j], extra = n%10, n/10
		}
		// fmt.Println(i, sumOfDigits(num), getPowerNum(num))
	}

	fmt.Println(topline, sumOfDigits(num), getPowerNum(num))
}

func sumOfDigits(num []byte) (s int) {
	var last int
	for last = len(num) - 1; num[last] == 0; last-- {
	}
	for i := 0; i <= last; i++ {
		s += int(num[i])
	}
	return
}

func getPowerNum(bytes []byte) string {
	num := make([]byte, len(bytes))
	copy(num, bytes)

	var last int
	for last = len(num) - 1; num[last] == 0; last-- {
	}
	for i, j := 0, last; i < j; i, j = i+1, j-1 {
		num[i], num[j] = num[j], num[i]
	}
	for i := 0; i <= last; i++ {
		num[i] += '0'
	}
	return string(num[:last+1])
}
