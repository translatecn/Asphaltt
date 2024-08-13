package main

import (
	"fmt"
	"strconv"
)

func main() {
	max := 0
	for x := 999; x >= 900; x-- {
		for y := 999; y >= 900; y-- {
			prod := x * y
			if prod > max && isPalindrome(prod) {
				max = prod
			}
		}
	}

	fmt.Println(max)
}

func isPalindrome(n int) bool {
	s := strconv.Itoa(n)
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}
	return true
}
