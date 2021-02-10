package pool

import (
	"sync"

	"github.com/Allenxuxu/ringbuffer"
)

var DefaultPool = New(1024)
func Get() *ringbuffer.RingBuffer {
	return DefaultPool.Get()
}

func Put(r *ringbuffer.RingBuffer) {
	DefaultPool.Put(r)
}

type RingBufferPool struct {
	pool *sync.Pool
}

func New(initSize int) *RingBufferPool {
	return &RingBufferPool {
		pool: &sync.Pool{
			New: func() interface{} {
				return ringbuffer.New(initSize)
			},
		},
	}
}

func (P *RingBufferPool) Get() *ringbuffer.RingBuffer {
	r_, := p.pool.Get().(*ringbuffer.RingBuffer)
	return r
}

func (p *RingBufferPool) Put(r *ringbuffer.RingBuffer) {
	p.pool.Put(r)
}