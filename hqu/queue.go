package hqu

import (
	"sync"
)

// Queue is a high memory efficient queue.
// It uses bucket to cache data and reuses the buckets.
type Queue struct {
	sync.Mutex

	// front is the location where a value should be dequeued
	// rear is the location where a new value should be enqueued
	front, rear int

	buckets  [][]interface{}
	freelist [][]interface{}

	// for reallocation of buckets
	orgBuckets [][]interface{}
}

var _ Queuer = &Queue{}

// Size gets the count of elements in queue
func (q *Queue) Size() int {
	q.Lock()
	size := q.rear - q.front
	q.Unlock()
	return size
}

// Enqueue pushes a value into the queue.
func (q *Queue) Enqueue(v interface{}) {
	q.Lock()
	q.enqueue(v)
	q.Unlock()
}

// Enqueue1 pushes a value into the queue and returns the size of the queue.
func (q *Queue) Enqueue1(v interface{}) (size int) {
	q.Lock()
	q.enqueue(v)
	size = q.rear - q.front
	q.Unlock()
	return
}

func (q *Queue) enqueue(v interface{}) {
	bp := q.rear / bucketSize

	var bkt []interface{}
	if bp == len(q.buckets) {
		if useFreelist && len(q.freelist) != 0 {
			// reuse bucket
			idx := len(q.freelist) - 1
			bkt = q.freelist[idx]
			q.freelist[idx] = nil
			q.freelist = q.freelist[:idx]
		} else {
			// create bucket
			bkt = make([]interface{}, bucketSize)
		}
		// realloc freelist automatically when necessary, TestSlice in slice_test.go
		q.buckets = append(q.buckets, bkt)
		q.orgBuckets = q.buckets
	} else {
		bkt = q.buckets[bp]
	}

	bkt[q.rear%bucketSize] = v
	q.rear++
}

// Dequeue pops a value from the queue with locking.
func (q *Queue) Dequeue() (v interface{}, ok bool) {
	q.Lock()
	v, ok = q.Dequeue0()
	q.Unlock()
	return
}

// Dequeue0 pops a value from the queue without locking.
func (q *Queue) Dequeue0() (v interface{}, ok bool) {
	if q.rear == q.front {
		return nil, false
	}

	bkt := q.buckets[0]
	v, ok = bkt[q.front], true
	bkt[q.front] = nil // free the value

	q.front++
	if q.front == bucketSize {
		q.buckets[0] = nil // free the bucket
		q.buckets = q.buckets[1:]

		// recude memory usage when no reallocation(append in Enqueue), TestSlice1 in slice_test.go
		if len(q.buckets)<<2 <= cap(q.orgBuckets) { // the usage is less than or equal to a quater of the capacity
			tmp := make([][]interface{}, cap(q.orgBuckets)/2)
			n := copy(tmp, q.buckets)
			q.buckets = tmp[:n]
			q.orgBuckets = tmp
		}

		// reuse bucket
		if useFreelist && len(q.freelist) < maxFreelist {
			q.freelist = append(q.freelist, bkt)
		}

		q.front = 0
		q.rear -= bucketSize
	}

	return
}

// Range does range all elements in queue with mutex lock,
// so you can't do `Enqueue` or `Dequeue` in `Range`.
func (q *Queue) Range(handle func(v interface{}) bool) {
	q.Lock()
	for i := q.front; i < q.rear; i++ {
		bp, fp := i/bucketSize, i%bucketSize
		if !handle(q.buckets[bp][fp]) {
			break
		}
	}
	q.Unlock()
}
