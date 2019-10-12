package shardedmap

import (
	"sync"
)

// Implementation: This is a sharded map so that the cost of locking is
// distributed with the data, instead of a single lock.
// The optimal number of shards will probably depend on the number of system
// cores but we provide a general default.
type Uint64Map struct {
	shardCount uint64 // Don't alter after creation, no mutex here
	mutexes    []sync.RWMutex
	maps       []map[uint64]interface{}
}

// NewUint64Map ...
func NewUint64Map(shardCount int) *Uint64Map {
	if shardCount <= 0 {
		shardCount = defaultShards
	}

	sm := &Uint64Map{
		shardCount: uint64(shardCount),
		mutexes:    make([]sync.RWMutex, shardCount),
		maps:       make([]map[uint64]interface{}, shardCount),
	}

	for i := range sm.maps {
		sm.maps[i] = make(map[uint64]interface{})
	}

	return sm
}

func (sm *Uint64Map) pickShard(key uint64) uint64 {
	// Assumes keys are well distributed. In the (rare?) case that they are
	// evenly separated, this could lead to a "hot" shard. In that case a
	// hashed picker would be better (TODO as an option)
	return key % sm.shardCount
}

// Store ...
func (sm *Uint64Map) Store(key uint64, value interface{}) {
	shard := sm.pickShard(key)
	sm.mutexes[shard].Lock()
	sm.maps[shard][key] = value
	sm.mutexes[shard].Unlock()
}

// Load ...
func (sm *Uint64Map) Load(key uint64) (interface{}, bool) {
	shard := sm.pickShard(key)
	sm.mutexes[shard].RLock()
	value, ok := sm.maps[shard][key]
	sm.mutexes[shard].RUnlock()
	return value, ok
}

// LoadOrStore ...
func (sm *Uint64Map) LoadOrStore(key uint64, value interface{}) (actual interface{}, loaded bool) {
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
func (sm *Uint64Map) Delete(key uint64) {
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
func (sm *Uint64Map) Range(f func(key uint64, value interface{}) bool) {
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
func (sm *Uint64Map) ConcRange(f func(key uint64, value interface{}) bool) {
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
func (sm *Uint64Map) AsyncRange(f func(key uint64, value interface{}) bool) {
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