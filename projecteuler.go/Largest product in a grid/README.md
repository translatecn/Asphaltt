# Largest product in a grid

Brute force to calculate all products in the grid.

## My solution

Do not forget to calculate the products in the last 3 rows and the last 3 columns.

```go
	nRow, nCol := len(numbers), len(numbers[0])
	maxProd := 0

	for i := 0; i < nRow-3; i++ {
		for j := 0; j < nCol-3; j++ {
			a := numbers[i][j] * numbers[i+1][j] * numbers[i+2][j] * numbers[i+3][j]
			b := numbers[i][j] * numbers[i][j+1] * numbers[i][j+2] * numbers[i][j+3]
			c := numbers[i][j] * numbers[i+1][j+1] * numbers[i+2][j+2] * numbers[i+3][j+3]
			d := numbers[i+3][j] * numbers[i+2][j+1] * numbers[i+1][j+2] * numbers[i][j+3]
			maxProd = max(a, b, c, d, maxProd)
		}
	}

	for i := 0; i < nRow-3; i++ {
		a := numbers[i][nCol-3] * numbers[i+1][nCol-3] * numbers[i+2][nCol-3] * numbers[i+3][nCol-3]
		b := numbers[i][nCol-2] * numbers[i+1][nCol-2] * numbers[i+2][nCol-2] * numbers[i+3][nCol-2]
		c := numbers[i][nCol-1] * numbers[i+1][nCol-1] * numbers[i+2][nCol-1] * numbers[i+3][nCol-1]
		maxProd = max(a, b, c, maxProd)
	}

	for j := 0; j < nCol-3; j++ {
		a := numbers[nRow-3][j] * numbers[nRow-3][j+1] * numbers[nRow-3][j+2] * numbers[nRow-3][j+3]
		b := numbers[nRow-2][j] * numbers[nRow-2][j+1] * numbers[nRow-2][j+2] * numbers[nRow-2][j+3]
		c := numbers[nRow-1][j] * numbers[nRow-1][j+1] * numbers[nRow-1][j+2] * numbers[nRow-1][j+3]
		maxProd = max(a, b, c, maxProd)
	}
```

[Next problem](https://github.com/Asphaltt/projecteuler.go/tree/main/Highly%20divisible%20triangular%20number).
