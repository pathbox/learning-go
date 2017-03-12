// Package cmap Concurrent safe Map based on map and sync.RWMutex
package cmap

import (
	"sync"
)

type Cmap struct {
	m map[interface{}]interface{}
	sync.RWMutex
}

// New new cmap
func New() *Cmap {
	return &Cmap{
		m: make(map[interface{}]interface{}),
	}
}

// Set maps the given key and value.
// Return false if the key is already in the map.
func (cm *Cmap) Set(k interface{}, v interface{}) bool {
	cm.Lock()
	defer cm.Unlock()
	_, ok := cm.m[k]
	if !ok { // is key is not already existing
		cm.m[k] = v
	}
	return !ok
}

func (cm *Cmap) Delete(k interface{}) {
	cm.Lock()
	defer cm.Unlock()
	delete(cm.m, k)
}

// Has returns true if k exists in the map

func (cm *Cmap) Has(k interface{}) bool {
	cm.RLock()
	defer cm.RUnlock()
	if _, ok := cm.m[k]; !ok {
		return false
	}
	return true
}

func (cm *Cmap) Len() int {
	cm.RLock()
	defer cm.RUnlock()
	return len(cm.m)
}

func (cm *Cmap) IsEmpty() bool {
	return cm.Len() == 0
}

func (cm *Cmap) Clear() {
	cm.Lock()
	defer cm.Unlock()
	cm.m = make(map[interface{}]interface{})
}

// Keys return all the keys in cmap
func (cm *Cmap) Keys() []interface{} {
	cm.RLock()
	defer cm.RUnlock()
	s := make([]interface{}, cm.Len())
	for k := range cm.m {
		s = append(s, k)
	}
	return s
}

// Values return all the values in cmap
func (cm *Cmap) Values() []interface{} {
	cm.RLock()
	defer cm.RUnlock()

	s := make([]interface{}, cm.Len())
	for _, v := range cm.m {
		s = append(s, v)
	}
	return s
}
