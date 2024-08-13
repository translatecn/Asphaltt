package main

import (
	"fmt"
)

func main() {
	const grid = 20
	lattices := make([][]int64, grid)
	for i := 0; i < grid; i++ {
		lattices[i] = make([]int64, grid)
	}

	lattices[grid-1][grid-1] = 2
	for i := grid - 2; i >= 0; i-- {
		lattices[i][grid-1] = lattices[i+1][grid-1] + 1
	}

	for j := grid - 2; j >= 0; j-- {
		lattices[grid-1][j] = lattices[grid-1][j+1] + 1
	}

	for i := grid - 2; i >= 0; i-- {
		for j := grid - 2; j >= 0; j-- {
			lattices[i][j] = lattices[i+1][j] + lattices[i][j+1]
		}
	}

	fmt.Println(lattices[0][0])
}
