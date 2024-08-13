# Summation of primes

Could you calculate all primes below two-million? No, please don't do that.

## My solution

Use algorithm of [the Sieve of Eratosthenes](https://en.wikipedia.org/wiki/Sieve_of_Eratosthenes).

Or you can use algorithm of [Euler's sieve](https://en.wikipedia.org/wiki/Sieve_of_Eratosthenes#Euler's_Sieve).

Here's the sieve of Eratosthenes.

```go
	const twoMillion = 2_000_000
	primes := make(map[int]struct{}) // here's a map instead of an array
	for i := 2; i < twoMillion; i++ {
		primes[i] = struct{}{}
	}

	for i := 2; i < twoMillion; i++ {
		if _, ok := primes[i]; !ok { // filter composite number
			continue
		}

		for j := 2; i*j <= twoMillion; j++ {
			delete(primes, i*j) // delete flags of one prime's multiples
		}
	}

	sum := int64(0)
	for k := range primes {
		sum += int64(k)
	}
```

[Next problem](https://github.com/Asphaltt/projecteuler.go/tree/main/Largest%20product%20in%20a%20grid).
