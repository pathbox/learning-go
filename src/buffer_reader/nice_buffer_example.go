package nicebuf

import (
	"io"
	"sync"
	"unsafe"
)

type Buffer struct {
	Bytes []byte
}

var _ io.Writer = &Buffer{}

// Write a chunk of bytes to the buffer
func (b *Buffer) Write(v []byte) (int, error) {
	b.Bytes = append(b.Bytes, v...) // nice way
	return len(v), nil
}

func (b *Buffer) WriteByte(v byte) {
	b.Bytes = append(b.Bytes, v)
}

func (b *Buffer) Reset() {
	b.Bytes = b.bytes[:0] // use this way to reset
}

func (b *Buffer) String() string {
	return *(*string)(unsafe.Pointer(&b.Bytes))
}

// WriteTo writes the contents of our buffer to an io.Writer
func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.Bytes)
	return int64(n), err
}

// Buffer pool
var bufpool = sync.Pool{
	New: func() interface{} { return &Buffer{} },
}

// NewBufferFromPool returns a pointer to a zerod Buffer. This may be retrieved from a
// pool. When you're done with it, call 'ReturnToPool'.
func NewBufferFromPool() *Buffer {
	b := bufpool.Get().(*Buffer)
	b.Reset()
	return b
}

// NewBufferFromPoolWithCap returns a pointer to a zero'd Buffer with its underlying
// capacity set. This may be retrieved from a pool. When you're done with it, call 'ReturnToPool'.
func NewBufferFromPoolWithCap(size int) *Buffer {
	b := bufpool.Get().(*Buffer)

	if c := cap(b.Bytes); c < size {
		b.Bytes = make([]byte, 0, size)
	} else if c > 0 {
		b.Reset()
	}

	return b
}

// ReturnToPool puts this instance back in the underlying pool. Reading from or using this instance
// in any way after calling this is invalid.
func (b *Buffer) ReturnToPool() {
	bufpool.Put(b)
}
