package shardedmap

import (
	"sync"
)

// Implementation: This is a sharded map so that the cost of locking is
// distributed with the data, instead of a single lock.
// The optimal number of shards will probably depend on the number of system
// cores but we provide a general default.

type StrMap struct {
	shardCount uint64 // Don't alter after creation, no mutex here
	mutexes    []sync.RWMutex
	maps       []map[string]interface{}
}

// NewStrMap ...
func NewStrMap(shardCount int) *StrMap {
	if shardCount <= 0 {
		shardCount = defaultShards
	}

	sm := &StrMap{
		shardCount: uint64(shardCount),
		mutexes:    make([]sync.RWMutex, shardCount),
		maps:       make([]map[string]interface{}, shardCount),
	}

	for i := range sm.maps {
		sm.maps[i] = make(map[string]interface{})
	}

	return sm
}

func (sm *StrMap) pickShard(key string) uint64 {
	return memHashString(key) % sm.shardCount
}

// Store ...
func (sm *StrMap) Store(key string, value interface{}) {
	shard := sm.pickShard(key)
	sm.mutexes[shard].Lock()
	sm.maps[shard][key] = value
	sm.mutexes[shard].Unlock()
}

// Load ...
func (sm *StrMap) Load(key string) (interface{}, bool) {
	shard := sm.pickShard(key)
	sm.mutexes[shard].RLock()
	value, ok := sm.maps[shard][key]
	sm.mutexes[shard].RUnlock()
	return value, ok
}

// LoadOrStore ...
func (sm *StrMap) LoadOrStore(key string, value interface{}) (actual interface{}, loaded bool) {
	shard := sm.pickShard(key)
	sm.mutexes[shard].RLock()
	// Fast path assuming value has a somewhat high chance of already being
	// there.
	if actual, loaded = sm.maps[shard][key]; loaded {
		sm.mutexes[shard].RUnlock()
		return
	}
	sm.mutexes[shard].RUnlock()
	// Gotta check again, unfortunately
	sm.mutexes[shard].Lock()
	if actual, loaded = sm.maps[shard][key]; loaded {
		sm.mutexes[shard].Unlock()
		return
	}
	sm.maps[shard][key] = value
	sm.mutexes[shard].Unlock()
	return value, loaded
}

// Delete ...
func (sm *StrMap) Delete(key string) {
	shard := sm.pickShard(key)
	sm.mutexes[shard].Lock()
	delete(sm.maps[shard], key)
	sm.mutexes[shard].Unlock()
}

// Range is modeled after sync.Map.Range. It calls f sequentially for each key
// and value present in each of the shards in the map. If f returns false, range
// stops the iteration.
//
// No key will be visited more than once, but if any value is inserted
// concurrently, Range may or may not visit it. Similarly, if a value is
// modified concurrently, Range may visit the previous or newest version of said
// value.
func (sm *StrMap) Range(f func(key string, value interface{}) bool) {
	for shard := range sm.mutexes {
		sm.mutexes[shard].RLock()
		for key, value := range sm.maps[shard] {
			if !f(key, value) {
				sm.mutexes[shard].RUnlock()
				return
			}
		}
		sm.mutexes[shard].RUnlock()
	}
}

// ConcRange ranges concurrently over all the shards, calling f sequentially
// over each shard's key and value. If f returns false, range stops the
// iteration on that shard (but the other shards continue until completion).
//
// No key will be visited more than once, but if any value is inserted
// concurrently, Range may or may not visit it. Similarly, if a value is
// modified concurrently, Range may visit the previous or newest version of said
// value.
func (sm *StrMap) ConcRange(f func(key string, value interface{}) bool) {
	var wg sync.WaitGroup
	wg.Add(int(sm.shardCount))
	for shard := range sm.mutexes {
		go func(shard int) {
			sm.mutexes[shard].RLock()
			for key, value := range sm.maps[shard] {
				if !f(key, value) {
					sm.mutexes[shard].RUnlock()
					wg.Done()
					return
				}
			}
			sm.mutexes[shard].RUnlock()
			wg.Done()
		}(shard)
	}
	wg.Wait()
}

// AsyncRange is exactly like ConcRange, but doesn't wait until all shards are
// done. This is usually ok, although calls that appear to happen "sequentially"
// on the same goroutine might get the before or after AsyncRange values, which
// might be surprising behaviour. When that's not desirable, use ConcRange.
func (sm *StrMap) AsyncRange(f func(key string, value interface{}) bool) {
	for shard := range sm.mutexes {
		go func(shard int) {
			sm.mutexes[shard].RLock()
			for key, value := range sm.maps[shard] {
				if !f(key, value) {
					sm.mutexes[shard].RUnlock()
					return
				}
			}
			sm.mutexes[shard].RUnlock()
		}(shard)
	}
}