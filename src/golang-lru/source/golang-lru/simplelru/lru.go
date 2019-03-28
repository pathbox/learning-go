package simplelru

import (
	"container/list"
	"errors"
)

type EvictCallback func(key interface{}, value interface{})

type LRU struct {
	size      int
	evictList *list.List                    // 1.列表存储底层数据
	items     map[interface{}]*list.Element // 2.map 从列表中同步数据，相当于上一层的缓存层，用于read操作，具有O(1)的读取性能，这样LRU Cache 就有高效的读性能
	onEvict   EvictCallback                 // 驱逐时的回调func
}

type entry struct {
	key   interface{}
	value interface{}
}

func NewLRU(size int, onEvict EvictCallback) (*LRU, error) {
	if size <= 0 {
		return nil, errors.New("Must provide a positive size")
	}
	c := &LRU{
		size:      size,
		evictList: list.New(),
		items:     make(map[interface{}]*list.Element),
		onEvict:   onEvict,
	}
	return c, nil
}

func (c *LRU) Purge() {
	for k, v := range c.items {
		if c.onEvict != nil {
			c.onEvict(k, v.Value.(*entry).value)
		}
		delete(c.items, k)
	}
	c.evictList.Init()
}

// Add adds a value to the cache.  Returns true if an eviction occurred.
func (c *LRU) Add(key, value interface{}) (evicted bool) {
	if ent, ok := c.items[key]; ok { // key原本存在，将其移到队列头部
		c.evictList.MoveToFront(ent)
		ent.Value.(*entry).value = value //修改过值
		return false
	}

	// Add new item
	ent := &entry{key, value}
	entry := c.evictList.PushFront(ent) // push到头部
	c.items[key] = entry

	evict := c.evictList.Len() > c.size
	if evict {
		c.removeOldest()
	}
	return evict
}

// Get looks up a key's value from the cache
func (c *LRU) Get(key interface{}) (value interface{}, ok bool) {
	if ent, ok := c.items[key]; ok { // 从map中取，如果不存在则返回nil和false
		c.evictList.MoveToFront(ent) // 并且将该ent 移到队列头部
		return ent.Value.(*entry).value, true
	}
	return
}

// Contains checks if a key is in the cache, without updating the recent-ness
// or deleting it for being stale.
func (c *LRU) Contains(key interface{}) (ok bool) {
	_, ok = c.items[key]
	return ok
}

// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *LRU) Peek(key interface{}) (value interface{}, ok bool) {
	var ent *list.Element
	if ent, ok = c.items[key]; ok {
		return ent.Value.(*entry).value, true
	}
	return nil, ok
}

func (c *LRU) Remove(key interface{}) (paresent bool) {
	if ent, ok := c.items[key]; ok {
		c.removeElement(ent)
		return true
	}
	return false
}

// RemoveOldest removes the oldest item from the cache.
func (c *LRU) RemoveOldest() (key interface{}, value interface{}, ok bool) {
	ent := c.evictList.Back()
	if ent != nil {
		c.removeElement(ent)
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return nil, nil, false
}

// GetOldest returns the oldest entry
func (c *LRU) GetOldest() (key interface{}, value interface{}, ok bool) {
	ent := c.evictList.Back()
	if ent != nil {
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return nil, nil, false

}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *LRU) Keys() []interface{} {
	keys := make([]interface{}, len(c.items))
	i := 0
	for ent := c.evictList.Back(); ent != nil; ent = ent.Prev() {
		keys[i] = ent.Value.(*entry).key
		i++
	}
	return keys
}

// Len returns the number of items in the cache.
func (c *LRU) Len() int {
	return c.evictList.Len()
}

func (c *LRU) removeOldest() {
	ent := c.evictList.Back()
	if ent != nil {
		c.removeElement(ent)
	}
}

func (c *LRU) removeElement(e *list.Element) {
	c.evictList.Remove(e)
	kv := e.Value.(*entry)
	delete(c.items, kv.key)
	if c.onEvict != nil {
		c.onEvict(kv.key, kv.value)
	}
}
