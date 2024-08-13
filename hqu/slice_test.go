package hqu

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {
	s := (make([]byte, 512))
	fmt.Printf("original ptr: %p\n", s)

	s0 := s[256:]
	fmt.Printf("half ptr: %p\n", s0)

	s0[0] = 0xff
	fmt.Printf("org array slice: %v\n", s[255:260])
	fmt.Printf("s0 array slice: %v\n", s0[:4])
	s0[0] = 0x00

	if s[256] != s0[0] {
		t.Fail()
	}

	var arr [256]byte
	s1 := append(s0, arr[:]...) // realloc
	fmt.Printf("after append ptr: %p\n", s1)

	s1[0] = 0xff
	fmt.Printf("org array slice: %v\n", s[255:260])
	fmt.Printf("s0 array slice: %v\n", s0[:4])

	if s[256] == s1[0] {
		t.Fail()
	}

	// conclusion: the appended slice uses a new under array when necessary, not the original under array
}

func TestSlice1(t *testing.T) {
	s := make([]byte, 512)
	org := s[:]
	fmt.Printf("original cap: %d, len: %d\n", cap(s), len(s))

	s = s[256:]
	fmt.Printf("half cap: %d, len: %d\n", cap(s), len(s))

	s = s[1:]
	fmt.Printf("minus one cap: %d, len: %d\n", cap(s), len(s))

	s[0] = 0xff // affected the under array

	fmt.Printf("org array slice: %v\n", org[255:260])

	if org[257] != s[0] {
		t.Fail()
	}

	// conclusion: the reduced slice s uses the original under array, not a new under array
}

func TestSlice2(t *testing.T) {
	s := make([]byte, 512)
	org := s[:]
	fmt.Printf("original ptr: %p, slice: %p\n", s, org)

	s0 := s[:256]
	fmt.Printf("half ptr: %p\n", s0)

	s1 := s[:128]
	fmt.Printf("quater ptr: %p\n", s1)

	s1[0] = 0xff // affected the under array
	fmt.Printf("org array slice: %v\n", org[:4])

	if org[0] != s1[0] {
		t.Fail()
	}

	// conclusion: the reduced slice s1 uses the original under array, not a new under array
}
