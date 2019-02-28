package lru

import "sync"

// 缩写：Least Recently Used（ 最近 最久 使用），时间维度 命中的先后顺序

// 原则：若数据在最近一段时间内都未使用（读取或更新），则以后使用几率也很低，应被淘汰

//使用链表：由于缓存读写删都是高频操作，考虑使用写删都为 O(1) 的链表，而非写删都为 O(N) 的数组。
// 使用双链表：选用删除操作为 O(1) 的双链表而非删除为 O(N) 的单链表。
// 维护额外哈希表：链表查找必须遍历 O(N) 读取，可在缓存中维护 map[key]*Node 的哈希表来实现O(1) 的链表查找。
type Node struct {
	key        string
	val        interface{}
	prev, next *Node
}

type List struct {
	head, tail *Node
	size       int
}

func (l *List) Prepend(node *Node) *Node {
	if l.head == nil {
		l.head = node
		l.tail = node
	} else {
		node.prev = nil
		node.next = l.head
		l.head.prev = node
		l.head = node // 先插入到head之前，然后head指针指导node节点，node节点为head
	}
	l.size++
	return node
}

func (l *List) Remove(node *Node) *Node {
	if node == nil {
		return nil
	}
	prev, next := node.prev, node.next // 先定位当前要remove 的node的 prev next
	if prev == nil {
		l.head = next // 删除头结点
	} else {
		prev.next = next
	}

	if next == nil {
		l.tail = prev // 删除尾结点
	} else {
		next.prev = prev
	}

	l.size--
	node.prev, node.next = nil, nil
	return node
}

func (l *List) MoveToHead(node *Node) *Node {
	if node == nil {
		return nil
	}
	n := l.Remove(node)
	return l.Prepend(n)
}

func (l *List) Tail() *Node {
	return l.tail
}

func (l *List) Size() int {
	return l.size
}

type LRUCache struct {
	capacity int              // 缓存空间大小
	items    map[string]*Node // 用于快速的取缓存数据
	list     *List            // 真正的存缓存数据 用链表方法删除和增加节点
	lock     *sync.RWMutex
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]*Node),
		list:     new(List),
		lock:     &sync.RWMutex{},
	}
}

func (c *LRUCache) Set(k string, v interface{}) {
	//命中 已在缓存中
	c.lock.RLock()
	node, ok := c.items[k]
	c.lock.RUnlock()
	if ok {
		node.val = v // 命中后更新值
		c.lock.Lock()
		c.items[k] = c.list.MoveToHead(node) // 将新缓存移到list的前面
		c.lock.Unlock()
		return
	}

	// 未命中 新的key
	node = &Node{key: k, val: v}
	if c.capacity == c.list.size {
		tail := c.list.Tail()
		c.lock.Lock()
		delete(c.items, tail.key) // k-v 数据存储与 node 中
		c.lock.Unlock()
		c.list.Remove(tail)
	}
	c.lock.Lock()
	c.items[k] = c.list.Prepend(node) // 更新node位置
	c.lock.Unlock()
}

func (c *LRUCache) Get(k string) interface{} {
	c.lock.RLock()
	node, ok := c.items[k]
	c.lock.RUnlock()
	if ok {
		c.lock.Lock()
		c.items[k] = c.list.MoveToHead(node)
		c.lock.Unlock()
		return node.val
	}
	return nil
}
