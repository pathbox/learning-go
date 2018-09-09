package shardmap

import (
	"sync"
)

var SHARD_COUNT uint = 1024

type SafeMap []*InternalMapShard

type InternalMapShard struct {
	sync.RWMutex
	items map[string]interface{} // internal storage map
}

func NewSafeMap() SafeMap {
	internalMaps := make(SafeMap, SHARD_COUNT)
	var ctr uint
	for ctr = 0; ctr < SHARD_COUNT; ctr++ {
		internalMaps[ctr] = &InternalMapShard{items: make(map[string]interface{})}
	}
	return internalMaps
}

func (safeMap SafeMap) GetShard(key string) *InternalMapShard {
	// array indexes are unit ( https://stackoverflow.com/a/16427832/2679770 )
	// fnv is a hash functio
	return safeMap[uint(fnv32(key))%uint(SHARD_COUNT)]
}

func (safeMap SafeMap) Set(key string, value interface{}) {
	// Get the shard to write to
	shard := safeMap.GetShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}
