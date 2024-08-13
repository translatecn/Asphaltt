# Sum square difference

Obviously, the square of sum is larger than the sum of squares.

## My solution

Brute force to calculate the square of sum and the sum of squares.

```go
    squareOfSum, sumOfSquares := 0, 0
	for i := 1; i <= 100; i++ {
		squareOfSum, sumOfSquares = squareOfSum+i, sumOfSquares+i*i
	}
    sq := int64(squareOfSum) * int64(squareOfSum)
    // diff := sq - int64(sumOfSquares)
```

[Next problem](https://github.com/Asphaltt/projecteuler.go/tree/main/10001st%20prime).
