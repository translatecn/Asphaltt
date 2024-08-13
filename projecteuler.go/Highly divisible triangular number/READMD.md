# Highly divisible triangular number

How to count divisors of a number fastly?

## My solution

```go
func main() {
	var divisors int
	triangle, i := 1, 2
	for ; divisors <= fiveHundred; triangle, i = triangle+i, i+1 {
		divisors = getDivisors(triangle)
	}
	triangle -= i + 1 // here's the answer
}

func getDivisors(n int) (res int) {
	for i, j := 1, n; i <= j; i, j = i+1, n/(i+1) { // count divisors fastly
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
```
