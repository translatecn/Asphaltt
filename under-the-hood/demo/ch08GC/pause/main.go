package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/trace"
	"time"
)

const (
	windowSize = 200000
	msgCount   = 1000000
)

type channel [windowSize][]byte

var (
	best    time.Duration = time.Second
	bestAt  time.Time
	worst   time.Duration
	worstAt time.Time

	start = time.Now()
)

func mkMessage(n int) []byte {
	m := make([]byte, 1024)
	for i := range m {
		m[i] = byte(n)
	}
	return m
}

func sendMessage(c *channel, highID int) {
	start := time.Now()
	(*c)[highID%windowSize] = mkMessage(highID)
	end := time.Now()
	elapsed := end.Sub(start)
	if elapsed > worst {
		worst = elapsed
		worstAt = end
	}
	if elapsed < best {
		best = elapsed
		bestAt = end
	}
}

func measure() {
	var c channel
	for i := 0; i < msgCount; i++ {
		sendMessage(&c, i)
	}
	fmt.Printf("Best push time %v at %v, worst push time: %v at %v. Wall clock: %v \n", best, bestAt.Sub(start), worst, worstAt.Sub(start), time.Since(start))
}

func main() {
	f, _ := os.Create("trace.out")
	defer f.Close()
	trace.Start(f)
	defer trace.Stop()

	for i := 0; i < 5; i++ {
		measure()
		worst = 0
		best = time.Second
		runtime.GC()
	}
}
