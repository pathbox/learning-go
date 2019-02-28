package lfu

import "sync"

// 数据结构
// 依旧使用双向链表实现高效写删操作，但 LFU 淘汰原则是 使用次数，数据节点在链表中的位置与之无关。可按使用次数划分 频率梯队，数据节点使用一次就挪到高频梯队。此外维护 minFreq 表示最低梯队，维护 2 个哈希表：

// map[freq]*List 各频率及其链表
// map[key]*Node 实现数据节点的 O(1) 读

type Node struct {
	key        string
	val        interface{}
	freq       int // 将节点从旧梯队移除时使用，非冗余存储
	prev, next *Node
}

type List struct {
	head, tail *Node
	size       int
}

type LFUCache struct {
	capacity int
	minFreq  int // 最低频率

	items map[string]*Node
	freqs map[int]*List // 不同频率梯队.每个梯队是一个List
	iLock *sync.RWMutex
}

func NewLFUCache(capacity int) *LFUCache {
	return &LFUCache{
		capacity: capacity,
		minFreq:  0,
		items:    make(map[string]*Node),
		freqs:    make(map[int]*List),
		iLock:    &sync.RWMutex{},
	}
}

func (c *LFUCache) Get(k string) interface{} {
	c.iLock.RLock()
	node, ok := c.items[k]
	c.iLock.RUnlock()
	if !ok {
		return -1
	}
	// 移到 +1 梯队中
	c.freqs[node.freq].Remove(node)
	node.freq++
	if _, ok := c.freqs[node.freq]; !ok {
		c.freqs[node.freq] = NewList()
	}
	newNode := c.freqs[node.freq].Prepend(node)
	c.items[k] = newNode // 新地址更新到 map
	if c.freqs[c.minFreq].Size() == 0 {
		c.minFreq++ // Get 的正好是当前值
	}
	return newNode.val
}

func (c *LFUCache) Set(k string, v interface{}) {
	if c.capacity <= 0 {
		return
	}

	// 命中，需要更新频率
	if val := c.Get(k); val != -1 {
		c.items[k].val = v // 直接更新值即可
		return
	}

	node := &Node{key: k, val: v, freq: 1}

	// 未命中
	// 缓存已满
	if c.capacity == len(c.items) {
		old := c.freqs[c.minFreq].Tail()
		c.freqs[c.minFreq].Remove(old)
		delete(c.items, old.key)
	}

	// 缓存未满，放入第 1 梯队
	c.items[k] = node
	if _, ok := c.freqs[1]; !ok {
		c.freqs[1] = NewList()
	}
	c.freqs[1].Prepend(node)
	c.minFreq = 1
}
