package main

import (
	"fmt"
)

func main() {
	// const twoMillion = 10
	const twoMillion = 2_000_000
	primes := make(map[int]struct{})
	for i := 2; i < twoMillion; i++ {
		primes[i] = struct{}{}
	}

	for i := 2; i < twoMillion; i++ {
		if _, ok := primes[i]; !ok {
			continue
		}

		for j := 2; j <= twoMillion/i; j++ {
			delete(primes, i*j)
		}
	}

	sum := int64(0)
	for k := range primes {
		sum += int64(k)
	}

	fmt.Println(sum)
}
