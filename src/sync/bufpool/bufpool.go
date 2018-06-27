package bufpool

import (
	"bytes"
	"sync"
)

// BufferPool represents a thread safe buffer bytes pool
type BufferPool struct {
	sync.Pool
}

func NewBufferPool(bufSize int) (bp *BufferPool) {
	return &BufferPool{
		sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, bufSize))
			},
		},
	}
}

// Get gets a Buffer from the SizedBufferPool, or creates a new one if none are
// available in the pool. Buffers have a pre-allocated capacity.
func (bp *BufferPool) Get() *bytes.Buffer {
	return bp.Pool.Get().(*bytes.Buffer)
}

// Put returns the given Buffer back to the SizedBufferPool
func (bp *BufferPool) Put(b *bytes.Buffer) {
	b.Reset()
	bp.Pool.Put(b)
}

// 减少重复 创建 分片 *bytes.Buffer的操作，虽然会多占用一定内存,但一定程度可以减轻GC压力
