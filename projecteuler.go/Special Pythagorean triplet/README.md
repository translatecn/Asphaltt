# Special Pythagorean triplet

Brute force to solve it even if it takes some time.

## My solution

To reduce solving time, don't calculate square-root of number.

We can use producting and adding to replace square-root calculating.

```go
	for a := 1; a < 1000; a++ {
		for b := 999; b > a; b-- {
			c := 1000 - a - b

			aa, bb, cc := a*a, b*b, c*c
			if aa+bb == cc || aa+cc == bb {
				fmt.Println(a * b * c) // we find it
				return
			}
		}
	}
```

[Next problem](https://github.com/Asphaltt/projecteuler.go/tree/main/Summation%20of%20primes).
