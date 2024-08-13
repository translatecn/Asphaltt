# Largest palindrome product

Brute force to solve it.

## My solution

Brute force to solve it with a little tricky.

Do not iterate the 1000 number twice. Iterate the last 100 number twice instead.

```go
	// iterate the last 100 number twice
	for x := 999; x >= 900; x-- {
		for y := 999; y >= 900; y-- {
			prod := x * y
			if prod > max && isPalindrome(prod) {
				max = prod
			}
		}
	}
```

[Next problem](https://github.com/Asphaltt/projecteuler.go/tree/main/Smallest%20multiple).
