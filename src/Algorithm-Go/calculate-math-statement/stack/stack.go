package stack

import "sync"

type Item string

type ItemStack struct {
	items []string
	lock  sync.RWMutex
}

func (s *ItemStack) New() *ItemStack {
	s.items = []string{}
	return s
}

func (s *ItemStack) Push(t string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.items = append(s.items, t)
}

func (s *ItemStack) Pop() string {
	s.lock.Lock()
	defer s.lock.Unlock()
	item := s.items[len(s.items)-1]
	s.items = s.items[0 : len(s.items)-1]
	return item
}

func (s *ItemStack) Top() string {
	return s.items[len(s.items)-1]
}

func (s *ItemStack) IsEmpty() bool {
	return len(s.items) == 0
}
