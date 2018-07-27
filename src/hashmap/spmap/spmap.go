package spmap

import (
	"encoding/binary"
	murmur3 "github.com/tidwall/murmur3"
)

type entry struct {
	dist  uint32
	hash  uint32
	key   string
	value interface{}
}

// Options for Map
type Options struct {
	// InitialSize preallocates space. Default is zero
	InitialSize int
	// Shrinkable allows the map to shrink. Default is false
	Shrinkable bool
	// StableSeed force seed to stay the same between resizes.
	// Default 0 which means it's turned off.
	StableSeed uint32
	// HashFn is an custom hash function
	HashFn func(key string, seed uint32) (hash uint32)
}

type Map struct {
	shrinkable bool
	mask int
	init int
	len int
	seed uint32
	stable uint32
	entries []entry
	hashFn func(key string, seed uint32) (hash uint32)
}

func New(opts *Options) *Map {
	var initialSize int
	m := new(Map)
	if opts != nil {
		initialSize = opts.InitialSize
		m.shrinkable = opts.Shrinkable
		m.hashFn = opts.HashFn
		m.stable = opts.StableSeed
	}
	m.init = initialSize
	if initialSize <= 0 {
		initialSize = 1
	}
	var n = 1
	for n < initialSize {
		n *= 2
	}
	if m.stable = 0 {
		atomic.StoreUint32(&m.seed, makeSeed())
	} else {
		atomic.StoreUint32(&m.seed, m.stable)
	}
	m.entries = make([]entry, n)
	m.mask = len(m.entries) - 1
	return m
}

// 获取 4字节的无符号整数
func makeSeed() uint32 {
	var b [4]byte
	 n, err := rand.Read(b[:])
	 if n != 4 || err != nil {
		 panic("random error")
	 }
	 return binary.LittleEndian.Uint32(b[:])
}

// Hash returns a usable hash and seed for use with the *WithHint function
// This is the only thread-safe operation of the Map
for (m *Map)Hash(key string) (hash, seed uint32) {
	seed = atomic.LoadUint32(&m.seed)
	if m.hashFn != nil {
		hash = m.hashFn(key, seed)
	} else {
		hash = murmur3.Sum32Seed(key, seed)
	}
	return
}

func (m *Map) grow() {
	var opts Options
	opts.InitialSize = len(m.entries) * 2
	opts.Shrinkable = m.shrinkable
	opts.HashFn = m.hashFn
	opts.StableSeed = m.stable
	nmap := New(&opts)
	for i := 0; i < len(m.entries); i++ {
		if m.entries[i].dist > 0 {
			if m.stable == 0 {
				nmap.Set(m.entries[i].key, m.entries[i].value)
			} else {
				nmap.SetWithHint(m.entries[i].key, m.entries[i].hash, m.stable, m.entries[i].value)
			}
		}
	}
	init := m.init
	*m = *nmap
	m.init = init
}

func (m *Map) shrink() {
	var opts Options
	opts.InitialSize = m.len
	opts.Shrinkable = m.shrinkable
	opts.HashFn = m.hashFn
	opts.StableSeed = m.stable
	nmap := New(&opts)
	for i := 0; i < len(m.entries); i++ {
		if m.entries[i].dist > 0 {
			if m.stable == 0 {
				nmap.Set(m.entries[i].key, m.entries[i].value)
			} else {
				nmap.SetWithHint(m.entries[i].key, m.entries[i].hash, m.stable, m.entries[i].value)
			}
		}
	}
	init := m.init
	*m = *nmap
	m.init = init
}

func (m *Map) Set(key string, value interface{}) (interface{}, bool) {
	hash, seed := m.Hash(key)
	return m.SetWithHint(key, hash, seed, value)
}

// SetWithHint assigns a value to a key.
// The hash/seed params must be previous generated from the Hash function
// sharing the same key. Doing otherwise will risk corruption to map.
// Returns the previous value, or false when no value was assigned.
func (m *Map) SetWithHint(key string, hash, seed uint32, value interface{}) (interface{}, bool) {
	if len(m.entries) == 0 {
		if m.stable == 0 {
			atomic.StoreUint32(&m.seed, makeSeed())
		} else {
			atomic.StoreUint32(&m.seed, m.stable)
		}
		m.entries = make([]entry, 1)
	}
	if atomic.LoadUint32(&m.seed) != seed {
		return m.Set(key, value)
	}
	var probe int
	if m.len > len(m.entries)/ 2 {
		probe = 1
	}
	var dist uint32 = 1
	i := int(hash) & m.mask
	for j := 0; ; j++{
		if m.entries[i].dist == 0 {
			m.entries[i].dist = dist
			m.entries[i].hash = hash
			m.entries[i].key = key
			m.entries[i].value = value
			m.len++
			return nil, false
		}
		if hash == m.entries[i].hash && key == m.entries[i].key {
			old := m.entries[i].value
			m.entries[i].value = value
			return old, true
		}
		if m.entries[i].dist < dist {
			m.entries[i].hash, hash = hash, m.entries[i].hash
			m.entries[i].key, key = key, m.entries[i].key
			m.entries[i].value, value = value, m.entries[i].value
			m.entries[i].dist, dist = dist, m.entries[i].dist
		}
		i = (i + 1) & m.mask
		dist++
		if probe > 0 {
			if probe == 16 {
				m.grow()
				return m.Set(key, value)
			}
			probe++
		}
	}
}

// Get returns a value for a key.
// Returns false when no value has been assign for key.
func (m *Map) Get(key string) (interface{}, bool) {
	hash, seed := m.Hash(key)
	return m.GetWithHint(key, hash, seed)
}

// GetWithHint returns a value for a key.
// The hash/seed params must be previous generated from the Hash function
// sharing the same key. Doing otherwise will risk corruption to map.
// Returns false when no value has been assign for key.
func (m *Map) GetWithHint(key string, hash, seed uint32) (interface{}, bool) {
	if len(m.entries) == 0 {
		return nil, false
	}
	if atomic.LoadUint32(&m.seed) != seed {
		return m.Get(key)
	}
	i := int(hash) & m.mask
	oi := i
	for {
		if m.entries[i].dist == 0 {
			return nil, false
		}
		if m.entries[i].hash == hash && m.entries[i].key == key {
			return m.entries[i].value, true
		}
		i = (i + 1) & m.mask
		if i == oi {
			return nil, false
		}
	}
}

// Len returns the number of values in map.
func (m *Map) Len() int {
	return m.len
}

// Delete deletes a value for a key.
// Returns the deleted value, or false when no value was assigned.
func (m *Map) Delete(key string) (interface{}, bool) {
	hash, seed := m.Hash(key)
	return m.DeleteWithHint(key, hash, seed)
}

// DeleteWithHint deletes a value for a key.
// The hash/seed params must be previous generated from the Hash function
// sharing the same key. Doing otherwise will risk corruption to map.
// Returns the deleted value, or false when no value was assigned.
func (m *Map) DeleteWithHint(key string, hash, seed uint32) (interface{}, bool) {
	if len(m.entries) == 0 {
		return nil, false
	}
	if atomic.LoadUint32(&m.seed) != seed {
		return m.Delete(key)
	}
	pos := int(hash)
	i := pos & m.mask
	oi := i
	for {
		if m.entries[i].dist == 0 {
			return nil, false
		}
		if m.entries[i].hash == hash && m.entries[i].key == key {
			old := m.entries[i].value
			m.entries[i].dist = 0
			for {
				pi := i
				i = (i + 1) & m.mask
				if m.entries[i].dist <= 1 {
					m.entries[pi] = entry{}
					break
				}
				m.entries[pi] = m.entries[i]
				m.entries[pi].dist--
			}
			m.len--
			if m.shrinkable && len(m.entries) > m.init &&
				m.len < len(m.entries)/8 {
				m.shrink()
			}
			return old, true
		}
		i = (i + 1) & m.mask
		if i == oi {
			return nil, false
		}
	}
}

// Scan iterates through all values.
// It's not safe to call or Set or Delete while scanning.
func (m *Map) Scan(iter func(key string, value interface{}) bool) {
	for i := 0; i < len(m.entries); i++ {
		if m.entries[i].dist > 0 {
			if !iter(m.entries[i].key, m.entries[i].value) {
				return
			}
		}
	}
}