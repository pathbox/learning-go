package bpool

type SizedBufferPool struct {
	c chan *bytes.Buffer
	a int
}

// SizedBufferPool creates a new BufferPool bounded to the given size.
// size defines the number of buffers to be retained in the pool and alloc sets
// the initial capacity of new buffers to minimize calls to make().

func NewSizedBufferPool(size int, alloc int) (bp *SizedBufferPool) {
	return &SizedBufferPool{
		c: make(chan *bytes.Buffer, size),
		a: alloc,
	}
}

// Get gets a Buffer from the SizedBufferPool, or creates a new one if none are
// available in the pool. Buffers have a pre-allocated capacity.

func (bp *SizedBufferPool) Get() (b *bytes.Buffer) {
	select {
	case b = <-bp.c:
		// return existing buffer
	default:
		// create new buffer
		b = bytes.NewBuffer(make([]byte, 0, bp.a))
	}
	return
}

// Put returns the given Buffer to the SizedBufferPool.
func (bp *SizedBufferPool) Put(b *bytes.Buffer) {
	b.Reset()
	// Release buffers over our maximum capacity and re-create a pre-sized
	// buffer to replace it.
	if cap(b.Bytes()) > bp.a {
		b = bytes.NewBuffer(make([]byte, 0, bp.a))
	}
	select {
	case bp.c <- b:
	default: // Discard the buffer if the pool is full.
	}
}
