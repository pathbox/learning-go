package hiset

import "encoding/json"

type set[T comparable] struct {
	setMap map[T]struct{}
}

func newSet[T comparable]() *set[T] {
	return &set[T]{
		setMap: make(map[T]struct{}),
	}
}

func (s *set[T]) Add(items ...T) bool {
	changed := false

	for index, item := range items {
		if _, ok := s.setMap[item]; ok {
			continue
		}
		changed = true
		s.setMap[items[index]] = struct{}{}
	}
	return changed
}

// Clear 清空数据
func (s *set[T]) Clear() {
	s.setMap = make(map[T]struct{})
}

// Size 返回set长度
func (s *set[T]) Size() int {
	return len(s.setMap)
}

func (s *set[T]) Contains(items ...T) bool {
	for _, item := range items {
		if _, ok := s.setMap[item]; ok {
			continue
		}
		return false
	}
	return true
}

// Equals 校验两个set是否相等
func (s *set[T]) Equals(cmpSet Set[T]) bool {
	if cmpSet == nil {
		return false
	}

	if cmpSet.Size() != len(s.setMap) {
		return false
	}

	for _, item := range cmpSet.ToArray() {
		if _, ok := s.setMap[item]; !ok {
			return false
		}
	}
	return true
}

// Union 取并集
func (s *set[T]) Union(other Set[T]) Set[T] {
	n := newSet[T]()
	union[T](s, other, n)
	return n
}

// Intersection 取交集
func (s *set[T]) Intersection(other Set[T]) Set[T] {
	n := newSet[T]()
	intersection[T](s, other, n)
	return n
}

// IsEmpty 校验set是否为空
func (s *set[T]) IsEmpty() bool {
	return len(s.setMap) == 0
}

// Remove 删除匹配的数据
func (s *set[T]) Remove(items ...T) bool {
	changed := false

	for _, item := range items {
		if _, ok := s.setMap[item]; !ok {
			continue
		}
		delete(s.setMap, item)
		changed = true
	}

	return changed
}

func (s *set[T]) ToArray() []T {
	var ret []T
	for key, _ := range s.setMap {
		ret = append(ret, key)
	}
	return ret
}

// MarshalJSON marshal json implement
func (s *set[T]) MarshalJSON() ([]byte, error) {
	arr := s.ToArray()
	return json.Marshal(arr)
}

// UnmarshalJSON unmarshal from json
func (s *set[T]) UnmarshalJSON(b []byte) error {
	var arr []T
	err := json.Unmarshal(b, &arr)
	if err != nil {
		return err
	}

	s.setMap = make(map[T]struct{})
	for index := range arr {
		s.setMap[arr[index]] = struct{}{}
	}

	return nil
}
