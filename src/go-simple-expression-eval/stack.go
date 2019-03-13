package main

type Stack struct {
	buf []interface{}
}

func NewStack() *Stack {
	return &Stack{
		buf: make([]interface{}, 0),
	}
}

func (s *Stack) Push(val interface{}) {
	s.buf = append(s.buf, val)
}

func (s *Stack) Pop() interface{} {
	len := len(s.buf)

	if len < 1 {
		panic("Attempting to Pop() from empty stack")
	}

	val := s.buf[len-1]

	s.buf = s.buf[0 : len-1]

	return val
}

func (s *Stack) Top() interface{} {
	return s.buf[len(s.buf)-1]
}

func (s *Stack) Len() int {
	return len(s.buf)
}
