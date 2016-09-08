package main

import "fmt"

type S struct{}

func (s *S) addr(){
  fmt.Printf("%p\n", s)
}

func main() {
  var a, b S
  a.addr()
  b.addr()
}
