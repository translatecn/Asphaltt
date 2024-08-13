package main

import (
	"fmt"
)

func main() {
	squareOfSum, sumOfSquares := 0, 0
	for i := 1; i <= 100; i++ {
		squareOfSum, sumOfSquares = squareOfSum+i, sumOfSquares+i*i
	}
	sq := int64(squareOfSum) * int64(squareOfSum)
	fmt.Println(sq - int64(sumOfSquares))
}
