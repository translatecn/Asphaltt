# largest product in a series

Brute force to calculate all products in the series.

# My solution

Ignore the duplicate work of calculating product of 12 adjacent numbers.

```go
	maxProd, n := int64(1), 13
	for i := 0; i < len(nums)-n; i++ {
		prod := int64(1)
		for j := i; j < i+n; j++ {
			prod *= int64(nums[j] - '0') // recalculate some numbers
		}
		if prod > maxProd {
			maxProd = prod
		}
	}
```

[Next problem](https://github.com/Asphaltt/projecteuler.go/tree/main/Special%20Pythagorean%20triplet).
