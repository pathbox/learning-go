package lrucache

// Element - node to store cache item
type Element struct {
	prev, next *Element    // 指针节点 链表
	Key        interface{} // interface
	Value      interface{}
}

// Next - fetch older element
func (e *Element) Next() *Element {
	return e.next
}

func (e *Element) Prev() *Element {
	return e.prev
}

// map + 双向链表
type LRUCache struct {
	cache    map[interface{}]*Element
	head     *Element
	tail     *Element
	capacity int
}

// New - create a new lru cache object
func New(capacity int) *LRUCache {
	return &LRUCache{make(map[interface{}]*Element, capacity), nil, nil, capacity}
}

// Put - put a cache item into lru cache
func (lc *LRUCache) Put(key interface{}, value interface{}) {
	if e, ok := lc.cache[key]; ok {
		e.Value = value
		lc.refresh(e)
		return
	}
	if lc.capacity == 0 {
		return
	} else if len(lc.cache) >= lc.capacity {
		// transfer the tail item as the new item, then refresh
		delete(lc.cache, lc.tail.Key)
		lc.tail.Key = key
		lc.tail.Value = value
		lc.cache[key] = lc.tail
		lc.refresh(lc.tail)
		return
	}

	e := &Element{nil, lc.head, key, value}
	lc.cache[key] = e
	if len(lc.cache) != 1 {
		lc.head.prev = e
	} else {
		lc.tail = e
	}
	lc.head = e // 将e节点放到头节点
}

// Get - get value of key from lru cache with result
func (lc *LRUCache) Get(key interface{}) (interface{}, bool) {
	if e, ok := lc.cache[key]; ok {
		lc.refresh(e)
		return e.Value, ok
	}
	return nil, false
}

// Delete - delete item by key from lru cache
func (lc *LRUCache) Delete(key interface{}) {
	if e, ok := lc.cache[key]; ok {
		delete(lc.cache, key)
		lc.remove(e)
	}
}

// Range - calls f sequentially for each key and value present in the lru cache
func (lc *LRUCache) Range(f func(key, value interface{}) bool) {
	for i := lc.head; i != nil; i = i.Next() {
		if !f(i.Key, i.Value) {
			break
		}
	}
}

// Update - inplace update
func (lc *LRUCache) Update(key interface{}, f func(value *interface{})) {
	if e, ok := lc.cache[key]; ok {
		f(&e.Value)
		lc.refresh(e)
	}
}

// Front - get front element of lru cache
func (lc *LRUCache) Front() *Element {
	return lc.head
}

// Back - get back element of lru cache
func (lc *LRUCache) Back() *Element {
	return lc.tail
}

// Len - length of lru cache
func (lc *LRUCache) Len() int {
	return len(lc.cache)
}

// Capacity - capacity of lru cache
func (lc *LRUCache) Capacity() int {
	return lc.capacity
}

func (lc *LRUCache) refresh(e *Element) {
	if e.prev != nil { // e不是head节点，将e修改为head节点  处理e的prev和next
		e.prev.next = e.next
		if e.next == nil {
			lc.tail = e.prev
		} else {
			e.next.prev = e.prev
		}
		e.prev = nil
		e.next = lc.head
		lc.head.prev = e
		lc.head = e
	}
}

func (lc *LRUCache) remove(e *Element) {
	// 处理e的prev和next
	// e是head情况
	if e.prev == nil {
		lc.head = e.next
	} else {
		e.prev.next = e.next
	}
	// e是tail情况
	if e.next == nil {
		lc.tail = e.prev
	} else {
		e.next.prev = e.prev
	}
}
