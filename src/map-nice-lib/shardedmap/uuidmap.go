package shardedmap

import (
	"sync"
)

type UUID [16]byte

// Implementation: This is a sharded map so that the cost of locking is
// distributed with the data, instead of a single lock.
// The optimal number of shards will probably depend on the number of system
// cores but we provide a general default.
type UUIDMap struct {
	shardCount uint64 // Don't alter after creation, no mutex here
	mutexes    []sync.RWMutex
	maps       []map[UUID]interface{}
}

// NewUUIDMap ...
func NewUUIDMap(shardCount int) *UUIDMap {
	if shardCount <= 0 {
		shardCount = defaultShards
	}

	sm := &UUIDMap{
		shardCount: uint64(shardCount),
		mutexes:    make([]sync.RWMutex, shardCount),
		maps:       make([]map[UUID]interface{}, shardCount),
	}

	for i := range sm.maps {
		sm.maps[i] = make(map[UUID]interface{})
	}

	return sm
}

func (sm *UUIDMap) pickShard(key UUID) uint64 {
	return memHash(key[:]) % sm.shardCount
}

// Store ...
func (sm *UUIDMap) Store(key UUID, value interface{}) {
	shard := sm.pickShard(key)
	sm.mutexes[shard].Lock()
	sm.maps[shard][key] = value
	sm.mutexes[shard].Unlock()
}

// Load ...
func (sm *UUIDMap) Load(key UUID) (interface{}, bool) {
	shard := sm.pickShard(key)
	sm.mutexes[shard].RLock()
	value, ok := sm.maps[shard][key]
	sm.mutexes[shard].RUnlock()
	return value, ok
}

// LoadOrStore ...
func (sm *UUIDMap) LoadOrStore(key UUID, value interface{}) (actual interface{}, loaded bool) {
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
func (sm *UUIDMap) Delete(key UUID) {
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
func (sm *UUIDMap) Range(f func(key UUID, value interface{}) bool) {
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
func (sm *UUIDMap) ConcRange(f func(key UUID, value interface{}) bool) {
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
// done. Since each shard is locked with an RWLock, it might be safe to use, but
// concurrent reads elsewhere might get the pre-range values, so don't use this
// one unless you don't care about that.
func (sm *UUIDMap) AsyncRange(f func(key UUID, value interface{}) bool) {
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