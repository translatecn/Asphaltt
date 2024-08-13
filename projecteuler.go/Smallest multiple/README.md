# Smallest multiple

I don't know how to calculate the smallest multiple of 1 to 20.

But I can count the factors for the smallest multiple of 1 to 20.

## My solution

The factors for the smallest multiple of 1 to 10 are `9 8 7 5`.

The factors for the smallest multiple of 1 to 20 are `19 17 16 13 11 9 7 5`.

```go
	// calculate product of the factors
	factors := []int{19, 17, 16, 13, 11, 9, 7, 5}
	for _, f := range factors {
		prod *= f
	}
```

[Next problem](https://github.com/Asphaltt/projecteuler.go/tree/main/Sum%20square%20difference).
