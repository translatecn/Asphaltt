# Multiples of 3 and 5

Use boolean flag to avoid adding multiples of 15 twice.

## My solution

Use a 1000 size boolean-array to indicate multiples of 3 and 5.

First step: make boolean flags.
```go
	// indicate whether a number is multiple of 3 or 5
	for i := 1; i < 1000; i++ {
		if i%3 == 0 || i%5 == 0 {
			integers[i] = true
		}
	}
```

Second step: calculate the sum.
```go
	// add up the numbers whose value is true
	for i := 1; i < 1000; i++ {
		if integers[i] {
			sum += i
		}
	}
```

That's all. Lets go to [next problem](https://github.com/Asphaltt/projecteuler.go/tree/main/Even%20Fibonacci%20numbers).
