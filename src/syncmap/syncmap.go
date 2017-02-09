package syncmap

import (
	"math/rand"
	"sync"
	"time"
)

const (
	defaultShardCount uint8 = 32
)

type syncMap struct {
	items map[string]interface{}
	sync.RWMutex
}

type SyncMap struct {
	shardCount uint8
	shards     []*syncMap
}

func New() *SyncMap {
	return NewWithShard(defaultShardCount)
}

func NewWithShard(shardCount uint8) *SyncMap {
	if !isPowerOfTwo(shardCount) {
		shardCount = defaultShardCount
	}
	m := new(SyncMap)
	m.shardCount = shardCount
	m.shards = make([]*syncMap, m.shardCount)
	for i, _ := range m.shards {
		m.shards[i] = &syncmap{items: make(map[string]interface{})}
	}
	return m
}

func (m *SyncMap) locate(key string) *syncMap {
	return m.shards[bkdrHash(key)&uint32((m.shardCount-1))]
}

func (m *SyncMap) Get(key string) (value interface{}, ok bool) {
	shard := m.locate(key)
	shard.RLock()
	value, ok = shard.items[key]
	shard.Unlock()
	return
}

func (m *SyncMap) Set(key string, value interface{}) {
	shard := m.locate(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

func (m *SyncMap) Delete(key string) {
	shard := m.locate(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

func (m *SyncMap) Pop() (string, interface{}) {
	if m.Size() == 0 {
		panic("syncmap: map is empty")
	}

	var (
		key   string
		value interface{}
		found = false
		n     = int(m.shardCount)
	)

	for !found {
		idx := rand.Intn(n)
		shard := m.shards[idx]
		shard.Lock()
		if len(shard.items) > 0 {
			found = true
			for key, value = range shard.items {
				break
			}
			delete(shard.items, key)
		}
		shard.Unlock()
	}
	return key, value
}

func (m *SyncMap) Size() int {
	size := 0
	for _, shard := range m.shards {
		shard.RLock()
		size += len(shard.items)
		shard.Unlock()
	}
	return size
}

func (m *SyncMap) Flush() int {
	size := 0
	for _, shard := range m.shards {
		shard.Lock()
		size += len(shard.items)
		shard.items = make(map[string]interface{})
		shard.Unlock()
	}
	return size
}

// Returns a channel from which each key in the map can be read
func (m *SyncMap) IterKeys() <-chan string {
	ch := make(chan string)
	go func() {
		for _, shard := range m.shards {
			shard.RLock()
			for key, _ := range shard.items {
				ch <- key
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

type Item struct {
	Key   string
	Value interface{}
}

// Return a channel from which each item (key:value pair) in the map can be read

func (m *SyncMap) IterItems() <-chan Item {
	ch := make(chan Item)
	go func() {
		for _, shard := range m.shards {
			shard.RLock()
			for key, value := range shard.items {
				ch <- Item{key, value}
			}
			shard.Unlock()
		}
		close(ch)
	}()
	return ch
}

const seed uint32 = 131

func bkdrHash(str string) uint32 {
	var h uint32
	for _, c := range str {
		h = h*seed + uint32(c)
	}
	return h
}

func isPowerOfTwo(x uint8) bool {
	return x != 0 && (x&(x-1) == 0)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
