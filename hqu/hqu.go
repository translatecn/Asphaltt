package hqu

const (
	bucketSize  = 8
	maxFreelist = 64
)

var useFreelist = true

type sizer interface {
	Size() int
}

type ranger interface {
	Range(func(interface{}) bool)
}

// Queuer interface for queue
type Queuer interface {
	sizer
	ranger
	Enqueue(v interface{})
	Dequeue() (v interface{}, ok bool)
}

// Stacker interface for stack
type Stacker interface {
	sizer
	ranger
	Push(v interface{})
	Pop() (v interface{}, ok bool)
}

type node struct {
	v interface{}

	prev, next *node
}

var (
	null = &node{nil, nil, nil}
)
