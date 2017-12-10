package pool

import "sync"

var DefaultPool = New(128, 65536, nil)

// Get pulls object whose generic size is at least of given size. It also
// returns a real size of x for further pass to Put(). It returns -1 as real
// size for nil x. Size >-1 does not mean that x is non-nil, so checks must be
// done.
//
// Note that size could be ceiled to the next power of two.
//
// Get is a wrapper around DefaultPool.Get().

func Get(size int) (interface{}, int) {
	return DefaultPool.Get(size)
}

func Put(x interface{}, size int) {
	DefaultPool.Put(x, size)
}

type Pool struct {
	pool map[int]*sync.Pool
	init func(int) interface{}
}

// New optionally specifies a function to generate
// a value when Get would otherwise return nil.
// It may not be changed concurrently with calls to Get.
func New(min, max int, init func(int) interface{}) *Pool {
	return &Pool{
		pool: MakePoolMap(min, max),
		init: init,
	}
}

// Get pulls object whose generic size is at least of given size. It also
// returns a real size of x for further pass to Put(). It returns -1 as real
// size for nil x. Size >-1 does not mean that x is non-nil, so checks must be
// done.
//
// Note that size could be ceiled to the next power of two.

func (p *Pool) Get(size int) (interface{}, int) {
	n := CeilToPowerOfTwo(size)
	pool, ok := p.pool[n]
	if ok {
		if x := pool.Get(); x != nil {
			return x, n
		}
	}
	if p.init == nil {
		// Nothing more to do.
		return nil, -1
	}
	if ok {
		// There is a pool for such size.
		// So init padded.
		return p.init(n), n
	}
	// There are no pool for such size.
	// So init with raw size.
	return p.init(size), size
}

// Put takes x and its size for future reuse.
func (p *Pool) Put(x interface{}, size int) {
	if pool, ok := p.pool[size]; ok {
		pool.Put(x)
	}
}
