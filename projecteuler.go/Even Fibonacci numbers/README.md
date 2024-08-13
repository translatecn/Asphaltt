# Even Fibonacci numbers

As we can see, sequence of Fibonacci number is `1 2 3 5 8 13 21 34 ...`.

Do you find the even-odd pattern of the sequence?

## My solution

I find the even-odd pattern of the sequence is that, excepting the first two Fibonacci number, there's an even number for every three Fibonacci number.

```go
	// add the very third Fibonacci number every time of generating three Fibonacci number
	sum, cnt := 2, 0
	x, y := 1, 2
	for y < 4_000_000 {
		x, y = y, x+y

		cnt++
		if cnt == 3 {
			sum, cnt = sum+y, 0
		}
	}
```

[Next problem](https://github.com/Asphaltt/projecteuler.go/tree/main/Largest%20prime%20factor).
