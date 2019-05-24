package lru

type entryU64 struct {
	key   uint64
	value interface{}
	prev  *entryU64
	next  *entryU64
}

// CacheU64 is a standard non-thread safe fixed-sized lru cache where the
// keys are uint64s and the values are interfaces.
type CacheU64 struct {
	entries map[uint64]*entryU64 // 每个map item是一个链表的head指针，也就是一个链表
	size    int
	onEvict func(key uint64, value interface{})
	head    *entryU64
	tail    *entryU64
}

// NewU64 returns a standard non-thread safe fixed-sized lru cache where the
// keys are uint64s and the values are interfaces.
// The size must be a positive number. The onEvict param is an optional
// callback function that fires when entries are forced to be evicted.
func NewU64(size int, onEvict func(key uint64, value interface{})) *CacheU64 {
	if size <= 0 {
		panic("invalid size")
	}
	return &CacheU64{
		entries: make(map[uint64]*entryU64, int(float64(size)*1.5)),
		size:    size,
		onEvict: onEvict,
	}
}

// Set a cache entry.
func (c *CacheU64) Set(key uint64, value interface{}) {
	e := c.entries[key]
	if e == nil {
		e = &entryU64{key: key, value: value}
		c.entries[key] = e
		if c.head == nil {
			c.head = e
			c.tail = e
		} else {
			c.head.prev = e // 从头部插入节点
			e.next = c.head
			c.head = e // 重置头部指针
		}
	} else {
		e.value = value
		c.promote(e)
	}
	if len(c.entries) > c.size {
		evicted := c.tail
		delete(c.entries, c.tail.key)
		c.tail = c.tail.prev
		c.tail.next = nil
		if c.onEvict != nil {
			c.onEvict(evicted.key, evicted.value)
		}
	}
}

func (c *CacheU64) promote(e *entryU64) { // 将该元素提升到头部
	if c.head != e {
		if c.tail == e {
			c.tail = c.tail.prev
			c.tail.next = nil
		} else {
			e.prev.next = e.next
			e.next.prev = e.prev
		}
		e.prev = nil
		e.next = c.head
		c.head.prev = e
		c.head = e
	}
}

// Len returns the number of entries in cache.
func (c *CacheU64) Len() int {
	return len(c.entries)
}

// Get an entry from cache.
func (c *CacheU64) Get(key uint64) interface{} {
	e := c.entries[key]
	if e == nil {
		return nil
	}
	c.promote(e)
	return e.value
}

func (c *CacheU64) Delete(key uint64) {
	e := c.entries[key]
	if e == nil {
		return
	}
	delete(c.entries, key)
	if len(c.entries) == 0 {
		c.head = nil
		c.tail = nil
	} else if e == c.head {
		c.head = c.head.next
		c.head.prev = nil
	} else if e == c.tail {
		c.tail = c.tail.prev
		c.tail.next = nil
	} else {
		e.prev.next = e.next
		e.next.prev = e.prev
	}
}
