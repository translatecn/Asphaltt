package hqu

import (
	"testing"
)

func TestStackPush(t *testing.T) {
	hq := &Stack{}

	hq.Push("test")
	v, ok := hq.Pop()
	if !ok || v != "test" {
		t.Fail()
	}
}

func TestStackPop(t *testing.T) {
	hq := &Stack{}

	v, ok := hq.Pop()
	if ok || v != nil {
		t.Fail()
	}
}

func BenchmarkStackWithFreelist(b *testing.B) {
	useFreelist = true
	hq := &Stack{}
	for i := 0; i < b.N; i++ {
		hq.Push(i)
	}

	for i := 0; i < b.N; i++ {
		hq.Pop()
	}
}

func BenchmarkStackWithoutFreelist(b *testing.B) {
	useFreelist = false
	hq := &Stack{}
	for i := 0; i < b.N; i++ {
		hq.Push(i)
	}

	for i := 0; i < b.N; i++ {
		hq.Pop()
	}
}

func TestStack0(t *testing.T) {
	s := &Stack{}

	for i := 0; i < bucketSize+1; i++ {
		s.Push(i)
	}

	if len(s.buckets) != 2 {
		t.Logf("length of buckets is not 2, is %d", len(s.buckets))
		t.Fail()
	}

	s.Pop()

	if len(s.buckets) != 1 {
		t.Logf("length of buckets is not 1, is %d", len(s.buckets))
		t.Fail()
	}

	if len(s.freelist) != 1 {
		t.Logf("length of freelist is not 1, is %d", len(s.freelist))
		t.Fail()
	}
}

func TestStack1(t *testing.T) {
	s := &Stack{}

	for i := 0; i < bucketSize*maxFreelist+1; i++ {
		s.Push(i)
	}

	if len(s.buckets) != maxFreelist+1 {
		t.Logf("length of buckets is not %d, is %d", maxFreelist+1, len(s.buckets))
		t.Fail()
	}

	for i := 0; i < bucketSize*(maxFreelist*3/4)+1; i++ {
		s.Pop()
	}

	if len(s.buckets) != maxFreelist/4 {
		t.Logf("length of buckets is not %d, is %d", maxFreelist/4, len(s.buckets))
		t.Fail()
	}

	if cap(s.orgBuckets) != maxFreelist/2 {
		t.Logf("capacity of under array is not %d, is %d", maxFreelist/2, cap(s.orgBuckets))
		t.Fail()
	}
}

func TestStackSize(t *testing.T) {
	s := &Stack{}

	for i := 0; i < bucketSize+1; i++ {
		s.Push(i)
	}

	if s.Size() != bucketSize+1 {
		t.Logf("size of stack is not %d, is %d", bucketSize+1, s.Size())
		t.Fail()
	}

	for i := 0; i < bucketSize; i++ {
		s.Pop()
	}

	if s.Size() != 1 {
		t.Logf("size of stack is not %d, is %d", 1, s.Size())
		t.Fail()
	}
}

func TestStackRange(t *testing.T) {
	s := &Stack{}

	for i := 0; i < bucketSize+1; i++ {
		s.Push(i)
	}

	idx := 0
	s.Range(func(v interface{}) bool {
		if v != idx {
			t.Fatalf("the `%d`th element is not %d, is %v", idx, idx, v)
			return false
		}
		idx++
		return true
	})

	for i := 0; i < bucketSize; i++ {
		s.Pop()
	}

	s.Range(func(v interface{}) bool {
		if v != 0 {
			t.Fatalf("the remain one element is not 0, is %v", v)
			return false
		}
		return true
	})
}
