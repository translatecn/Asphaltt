package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"math"
)

var (
	flgLarge    = flag.Int("L", 0, "only primes less than or equal to the number will be generated")
	flgJSONFile = flag.String("o", "primes.json", "save primes to the file")
)

func main() {
	flag.Parse()
	N := *flgLarge
	var x, y, n int
	nsqrt := math.Sqrt(float64(N))

	isPrimes := make([]bool, *flgLarge)

	for x = 1; float64(x) <= nsqrt; x++ {
		for y = 1; float64(y) <= nsqrt; y++ {
			n = 4*(x*x) + y*y
			if n <= N && (n%12 == 1 || n%12 == 5) {
				isPrimes[n] = !isPrimes[n]
			}
			n = 3*(x*x) + y*y
			if n <= N && n%12 == 7 {
				isPrimes[n] = !isPrimes[n]
			}
			n = 3*(x*x) - y*y
			if x > y && n <= N && n%12 == 11 {
				isPrimes[n] = !isPrimes[n]
			}
		}
	}

	for n = 5; float64(n) <= nsqrt; n++ {
		if isPrimes[n] {
			for y = n * n; y < N; y += n * n {
				isPrimes[y] = false
			}
		}
	}

	isPrimes[2] = true
	isPrimes[3] = true

	primes := make([]int, 0, 1270606)
	for x = 0; x < len(isPrimes)-1; x++ {
		if isPrimes[x] {
			primes = append(primes, x)
		}
	}

	// primes is now a slice that contains all the
	// primes numbers up to N

	data, _ := json.Marshal(primes)
	ioutil.WriteFile(*flgJSONFile, data, 0644)
}
