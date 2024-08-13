package hqu

import (
	"sync"
)

// Stack is a high memory efficient stack.
// It uses bucket to cache data and reuses the buckets.
type Stack struct {
	sync.Mutex

	pos int // the position where a new value should be inserted

	buckets  [][]interface{}
	freelist [][]interface{}

	// for reallocation of buckets
	orgBuckets [][]interface{}
}

var _ Stacker = &Stack{}

// Size gets the count of elements in stack
func (s *Stack) Size() int {
	s.Lock()
	size := s.pos
	s.Unlock()
	return size
}

// Push pushes a value into the stack.
func (s *Stack) Push(v interface{}) {
	s.Lock()

	bp := s.pos / bucketSize

	// look up bucket
	var bkt []interface{}
	if bp == len(s.buckets) {
		if useFreelist && len(s.freelist) != 0 {
			// reuse bucket
			idx := len(s.freelist) - 1
			bkt = s.freelist[idx]
			s.freelist[idx] = nil
			s.freelist = s.freelist[:idx]
		} else {
			bkt = make([]interface{}, bucketSize) // create bucket
		}
		// realloc freelist automatically when necessary, TestSlice in slice_test.go
		s.buckets = append(s.buckets, bkt)
		s.orgBuckets = s.buckets
	} else {
		bkt = s.buckets[bp]
	}
	bkt[s.pos%bucketSize] = v
	s.pos++

	s.Unlock()
}

// Pop pops a value from the stack.
func (s *Stack) Pop() (v interface{}, ok bool) {
	s.Lock()
	if s.pos == 0 {
		s.Unlock()
		return nil, false
	}
	s.pos--

	bp, qp := s.pos/bucketSize, s.pos%bucketSize

	// lookup bucket
	bkt := s.buckets[bp]
	v, ok = bkt[qp], true
	bkt[qp] = nil // free the value

	if qp == 0 {
		s.buckets[bp] = nil // free the bucket
		s.buckets = s.buckets[:len(s.buckets)-1]

		// reduce memory usage when no reallocation(append in Push), TestSlice2 in slice_test.go
		if len(s.buckets)<<2 <= cap(s.orgBuckets) { // the usage is less than or equal to a quater of the capacity
			tmp := make([][]interface{}, cap(s.orgBuckets)/2)
			n := copy(tmp, s.buckets)
			s.buckets = tmp[:n]
			s.orgBuckets = tmp
		}

		// reuse bucket
		if useFreelist && len(s.freelist) < maxFreelist {
			s.freelist = append(s.freelist, bkt)
		}
	}

	s.Unlock()
	return
}

// Range does range all elements in stack with mutex lock,
// so you can't do `Push` or `Pop` in `Range`.
func (s *Stack) Range(handle func(v interface{}) bool) {
	s.Lock()
	for i := 0; i < s.pos; i++ {
		bp, qp := i/bucketSize, i%bucketSize
		if !handle(s.buckets[bp][qp]) {
			break
		}
	}
	s.Unlock()
}
