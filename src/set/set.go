
package safeset
import (
  "sync"
)

type Set struct {
  m map[interface{}] bool
  sync.RWMutex
}

func New() *Set {
  return &Set{
    m: make(map[interface{}]bool)
  }
}

func (s *Set) Add(item interface{}) {
  s.Lock()
  defer s.Unlock()
  s.m[item] = true
}

func (s *Set) Remove(item interface{}) {
  s.Lock()
  defer s.Unlock()
  delete(s.m, item)
}

func (s *Set) Has(item interface{}) bool {
  s.RLock()
  defer s.RUnlock()
  _, ok := s.m[item]
  return ok
}

func (s *Set) Clear() {
  s.Lock()
  defer s.Unlock()
  s.m = make(map[interface{}]bool)
}

func (s *Set) IsEmpty() bool {
  if s.Len() == 0 {
    return true
  }
  return false
}

func (s *Set) List() []interface{} {
  s.RLock()
  defer s.RUnlock()
  list := make([]interface{}, 0)
  for item := range s.m {
    list = append(list, item)
  }
  return list
}