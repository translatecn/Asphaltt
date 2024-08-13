# Largest prime factor

How to generate prime numbers?

# My solution

I use an [online tool](https://onlinemathtools.com/generate-prime-numbers) to generate prime numbers.

And I guest the answer is one of the very first prime numbers.

```go
	number := int64(600851475143)
	primes := []int64{2, 3, 5, 7, 11, 13, 17, ...} // the first 1000 prime numbers.
	maxFactor := int64(0)
	for _, n := range primes {
		for number != 0 && number%n == 0 {
			number, maxFactor = number/n, n
		}
		if number == 0 {
			break
		}
	}
```

[Next problem](https://github.com/Asphaltt/projecteuler.go/tree/main/Largest%20palindrome%20product).
