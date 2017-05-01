package main

import (
	"log"
)

type Session struct{}

func main() {
	s := new(Session)
	s.log("Hello World!")
	num := 100
	list := []int{1, 2, 3}
	s.log(num)
	s.log(list)
}

func (s *Session) log(args ...interface{}) {
	log.Println(args...)
}
