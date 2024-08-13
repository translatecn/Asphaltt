package hqu

import (
	"math"
)

// FiboN gets the `n`th Fibonacci Number which starts with 1, 1, 2, ...
// with iteration.
func FiboN(n int) int {
	if n <= 2 {
		return 1
	}

	x, y := 1, 1
	for i := 3; i <= n; i++ {
		x, y = y, x+y
	}
	return y
}

// FiboNRecursive gets the `n`th Fibonacci Number which starts with 1, 1, 2, ...
// with recursion.
func FiboNRecursive(n int) int {
	if n <= 2 {
		return 1
	}
	return fiboNRecursive(n-2, 1, 1)
}

func fiboNRecursive(n, x, y int) int {
	if n == 0 {
		return y
	}
	return fiboNRecursive(n-1, y, x+y)
}

// FiboNPolynomial gets the `n`th Fibonacci Number which starts with 1, 1, 2, ...
// with polynomial formula.
func FiboNPolynomial(n int) int {
	if n < 3 {
		return 1
	}

	arg := math.Sqrt(5)
	return int(math.Round(1/arg*(math.Pow((1+arg)/2, float64(n))) + math.Pow((1-arg)/2, float64(n))))
	// return int(math.Round(1 / arg * ((pow((1+arg)/2, (n))) + (pow((1-arg)/2, (n)))))) // failed to compute the Fib result
}

func pow(x float64, n int) float64 { // n is positive
	if n == 0 {
		return 1.0
	}
	if n == 1 || x == 1.0 {
		return x
	}
	return pow0(x, 1.0, n)
}

func pow0(x, factor float64, n int) float64 {
	if n == 1 {
		return x * factor // complement a factor
	}
	return pow0(x*x*factor, pow(x, n%2), n/2)
}

// 0 -> 1; 1 -> x
// √ pow(x, 0/1)
// or 0 -> 1, x -> x
// or x -> 1, x+x -> x
// how to use '+-*/%', any number and 'x' to achieve it
