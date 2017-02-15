package main

import (
	"fmt"
	"regexp"
)

type I interface {
	Find(b []byte) []byte // re 有 Find(b []byte) 方法
}

func find(i I) {
	fmt.Printf("%s\n", i.Find([]byte("abc")))
}

func main() {
	var re = regexp.MustCompile(`b`)
	find(re)
}
