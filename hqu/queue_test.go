package hqu

import "testing"

func TestQueueEnqueue(t *testing.T) {
	hq := &Queue{}

	hq.Enqueue(1)

	v, ok := hq.Dequeue()
	if !ok || v != 1 {
		t.Fail()
	}
}

func TestQueueDequeue(t *testing.T) {
	hq := &Queue{}

	v, ok := hq.Dequeue()
	if ok || v != nil {
		t.Fail()
	}
}

func TestQueue0(t *testing.T) {
	q := &Queue{}

	for i := 0; i < 9; i++ {
		q.Enqueue(i)
	}

	if len(q.buckets) != 2 {
		t.Log("length of buckets is not 2")
		t.Fail()
	}

	for i := 0; i < 8; i++ {
		q.Dequeue()
	}

	if len(q.buckets) != 1 {
		t.Log("length of buckets is not 1")
		t.Fail()
	}

	if len(q.freelist) != 1 {
		t.Log("length of freelist is not 1")
		t.Fail()
	}
}

func TestQueue1(t *testing.T) {
	q := &Queue{}

	for i := 0; i < bucketSize*maxFreelist; i++ {
		q.Enqueue(i)
	}

	if q.Size() != bucketSize*maxFreelist {
		t.Logf("queue size, got %d, expect %d", q.Size(), bucketSize*maxFreelist)
		t.Fail()
	}

	if len(q.buckets) != maxFreelist {
		t.Logf("length of buckets is not %d, is %d", maxFreelist, len(q.buckets))
		t.Fail()
	}

	for i := 0; i < ((maxFreelist*3)/4)*bucketSize; i++ {
		q.Dequeue()
	}

	if len(q.buckets) != maxFreelist/4 {
		t.Logf("length of buckets is not %d, is %d", maxFreelist/4, len(q.buckets))
		t.Fail()
	}

	if cap(q.orgBuckets) != maxFreelist/2 {
		t.Logf("capacity of under array is not %d, is %d", maxFreelist/2, cap(q.orgBuckets))
		t.Fail()
	}
}

func TestQueueSize(t *testing.T) {
	q := &Queue{}

	for i := 0; i < bucketSize+1; i++ {
		q.Enqueue(i)
	}

	if q.Size() != bucketSize+1 {
		t.Logf("size of queue is not %d, is %d", bucketSize+1, q.Size())
		t.Fail()
	}

	for i := 0; i < bucketSize; i++ {
		q.Dequeue()
	}

	if q.Size() != 1 {
		t.Logf("size of queue is not %d, is %d", 1, q.Size())
		t.Fail()
	}
}

func TestQueueRange(t *testing.T) {
	q := &Queue{}

	for i := 0; i < bucketSize+1; i++ {
		q.Enqueue(i)
	}

	idx := 0
	q.Range(func(v interface{}) bool {
		if v != idx {
			t.Fatalf("the `%d`th element is not %d, is %v", idx, idx, v)
			return false
		}
		idx++
		return true
	})

	for i := 0; i < bucketSize; i++ {
		q.Dequeue()
	}

	q.Range(func(v interface{}) bool {
		if v != bucketSize {
			t.Fatalf("the remain one element is not %d, is %v", bucketSize, v)
			return false
		}
		return true
	})
}
