package bitcask

import (
	"errors"
	"sync"
)

type index struct {
	entries map[string]*entry
	*sync.RWMutex
}

var (
	ErrKeyNotFound = errors.New("Key not found")
)

func newIndex() *index {
	return &index{
		entries: make(map[string]*entry),
		RWMutex: &sync.RWMutex{},
	}
}

func (i *index) put(key string, entry *entry) {
	i.Lock()
	defer i.Unlock()
	i.entries[key] = entry
}

func (i *index) get(key []byte) (*entry, error) {
	i.Lock()
	defer i.Unlock()
	if entry, ok := i.entries[string(key)]; ok {
		return entry, nil
	}

	return nil, ErrKeyNotFound
}

func (i *index) del(key string) {
	i.Lock()
	defer i.Unlock()
	delete(i.entries, key)
}
