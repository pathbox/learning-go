package cache

import (
	"container/list"
	"sync"
)

// Cache is a bounded-size in-memory cache of sized items with a configurable eviction policy
type Cache interface {
	// Get retrieves items from the cache by key.
	// If an item for a particular key is not found, its position in the result will be nil.
	Get(keys ...string) []Item

	// Put adds an item to the cache.
	Put(key string, item Item)

	// Remove clears items with the given keys from the cache
	Remove(keys ...string)

	// Size returns the size of all items currently in the cache.
	Size() uint64
}

type Item interface {
	Size() uint64
}

// A tuple tracking a cached item and a reference to its node in the eviction list
type cached struct {
	item    Item
	element *list.Element
}

func (c *cached) setElementIfNotNil(element *list.Element) {
	if element != nil {
		c.element = element
	}
}

type cache struct {
	sync.Mutex                                  // Lock for synchronizing Get, Put, Remove
	cap          uint64                         // Capacity bound
	size         uint64                         // Cumulative size
	items        map[string]*cached             // Map from keys to cached items
	keyList      *list.List                     // List of cached items in order of increasing evictability
	recordAdd    func(key string) *list.Element // Function called to indicate that an item with the given key was added
	recordAccess func(key string) *list.Element // Function called to indicate that an item with the given key was accessed
}

// CacheOption configures a cache.
type CacheOption func(*cache)

// Policy is a cache eviction policy for use with the EvictionPolicy CacheOption.
type Policy uint8

const (
	// LeastRecentlyAdded indicates a least-recently-added eviction policy.
	LeastRecentlyAdded Policy = iota // 0
	// LeastRecentlyUsed indicates a least-recently-used eviction policy.
	LeastRecentlyUsed // 1
)

// EvictionPolicy sets the eviction policy to be used to make room for new items.
// If not provided, default is LeastRecentlyUsed.
func EvictionPolicy(policy Policy) CacheOption {
	return func(c *cache) {
		switch policy {
		case LeastRecentlyAdded:
			c.recordAccess = c.noop
			c.recordAdd = c.record
		case LeastRecentlyUsed:
			c.recordAccess = c.record
			c.recordAdd = c.noop
		}
	}
}

// New returns a cache with the requested options configured.
// The cache consumes memory bounded by a fixed capacity,
// plus tracking overhead linear in the number of items.
func New(capacity uint64, options ...CacheOption) Cache {
	c := &cache{
		cap:     capacity,
		keyList: list.New(),
		items:   map[string]*cached{}, // 真正的数据存在cached中
	}
	// Default LRU eviction policy
	EvictionPolicy(LeastRecentlyUsed)(c)

	for _, option := range options {
		option(c)
	}

	return c
}

func (c *cached) Get(keys ...string) []Item {
	c.Lock()
	defer c.Unlock()

	items := make([]Item, len(keys))

	for i, key := range keys {
		cached := c.items[key]
		if cached == nil {
			items[i] = nil
		} else {
			c.recordAccess(key)
			items[i] = cahced.item
		}
	}
	return items
}

func (c *cache) Put(key string, item Item) {
	c.Lock()
	defer c.Unlock()

	c.remove(key)

	// Make sure there's room to add this item
	c.ensureCapacity(item.Size)

	cached := &cahced{item: item}
	cached.setElementIfNotNil(c.recordAdd(key))
	cached.setElementIfNotNil(c.recordAccess(key))
	c.items[key] = cached
	c.size += item.Size()

}

func (c *cached) Remove(keys ...string) {
	c.Lock()
	defer c.Unlock()
	for _, key := range keys {
		c.remove(key)
	}
}

func (c *cache) Size() uint64 {
	return c.size
}

// Remove the item associated with the given key.
// The caller should hold the cache lock.
func (c *cache) remove(key string) {
	if cached, ok := c.items[key]; ok {
		delete(c.items, key)
		c.size -= cached.item.Size()
		c.keyList.Remove(cached.element)
	}
}

// A no-op function that does nothing for the provided key
func (c *cache) noop(string) *list.Element { return nil }

// A function to record the given key and mark it as last to be evicted
func (c *cache) record(key string) *list.Element {
	if item, ok := c.items[key]; ok {
		c.keyList.MoveToFront(item.element)
		return item.element
	}
	return c.keyList.PushFront(key)
}
