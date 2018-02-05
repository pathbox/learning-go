/*
Package queue provides a fast, ring-buffer queue based on the version suggested by Dariusz Górecki.
Using this instead of other, simpler, queue implementations (slice+append or linked list) provides
substantial memory and time benefits, and fewer GC pauses.
The queue implemented here is as fast as it is for an additional reason: it is *not* thread-safe.
*/

package queue

const minQueueLen = 16

type Queue struct {
	buf               []interface{}
	head, tail, count int
}

func New() *Queue {
	return &Queue{
		buf: make([]interface{}, minQueueLen),
	}
}

func (q *Queue) Length() int {
	return q.count
}

// resizes the queue to fit exactly twice its current contents
// this can result in shrinking if the queue is less than half-full
func (q *Queue) resize() {
	newBuf := make([]interface{}, q.count<<1)

	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail]) // 将数据复制到newBuf
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}

// Add puts an element on the end of the queue
// 从尾部增加，操作tail索引
func (q *Queue) Add(elem interface{}) {
	if q.count == len(q.buf) {
		q.resize()
	}

	q.buf[q.tail] = elem
	q.tail = (q.tail + 1) & (len(q.buf) - 1)
	q.count++
}

func (q *Queue) Peek() interface{} {
	if q.count <= 0 {
		panic("queue: Peek() called on empty queue")
	}
	return q.buf[q.head]
}

// Get returns the element at index i in the queue. If the index is
// invalid, the call will panic. This method accepts both positive and
// negative index values. Index 0 refers to the first element, and
// index -1 refers to the last.
func (q *Queue) Get(i int) interface{} {
	if i < 0 {
		i += q.count
	}
	if i < 0 || i >= q.count {
		panic("queue: Get() called with index out of range")
	}

	return q.buf[(q.head+i)&(len(q.buf)-1)]
}

// 从头部取出，所以 操作head索引
func (q *Queue) Remove() interface{} {
	if q.count <= 0 {
		panic("queue: Remove() called on empty queue")
	}

	ret := q.buf[q.head]
	q.buf[q.head] = nil

	q.head = (q.head + 1) & (len(q.buf) - 1) // 0 or q.head+1
	q.count--
	if len(q.buf) > minQueueLen && (q.count<<2) == len(q.buf) {
		q.resize()
	}
	return ret
}

// (q.head + 1) & (len(q.buf)- 1) 相同时，得到 这个相同的数，不同时 得到0
