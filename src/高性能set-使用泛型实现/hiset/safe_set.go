package hiset

import (
	"sync"
)

// threadSafeSet 线程安全set
type threadSafeSet[T comparable] struct {
	baseSet set[T]
	mutex   sync.RWMutex
}

// newThreadSafeSet ...
func newThreadSafeSet[T comparable]() *threadSafeSet[T] {
	return &threadSafeSet[T]{
		baseSet: set[T]{
			setMap: make(map[T]struct{}),
		},
	}
}

// Add 添加数据
func (s *threadSafeSet[T]) Add(items ...T) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.baseSet.Add(items...)
}

// Clear 清空数据
func (s *threadSafeSet[T]) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.baseSet.Clear()
}

// Size 返回set长度
func (s *threadSafeSet[T]) Size() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.baseSet.Size()
}

// Contains 是否包含对应数据
func (s *threadSafeSet[T]) Contains(items ...T) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.baseSet.Contains(items...)
}

// Equals 校验两个set是否相等
func (s *threadSafeSet[T]) Equals(cmpSet Set[T]) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.baseSet.Equals(cmpSet)
}

// IsEmpty 校验set是否为空
func (s *threadSafeSet[T]) IsEmpty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.baseSet.IsEmpty()
}

// Remove 删除匹配的数据
func (s *threadSafeSet[T]) Remove(items ...T) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.baseSet.Remove(items...)
}

// ToArray 返回set列表
func (s *threadSafeSet[T]) ToArray() []T {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.baseSet.ToArray()
}

// MarshalJSON marshal json implement
func (s *threadSafeSet[T]) MarshalJSON() ([]byte, error) {
	return s.baseSet.MarshalJSON()
}

// UnmarshalJSON unmarshal from json
func (s *threadSafeSet[T]) UnmarshalJSON(b []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.baseSet.UnmarshalJSON(b)
}

// Union 取并集
func (s *threadSafeSet[T]) Union(other Set[T]) Set[T] {
	n := newThreadSafeSet[T]()
	union[T](s, other, n)
	return n
}

// Intersection 取交集
func (s *threadSafeSet[T]) Intersection(other Set[T]) Set[T] {
	n := newThreadSafeSet[T]()
	intersection[T](s, other, n)
	return n
}
