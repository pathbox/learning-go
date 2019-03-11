package main

import (
	"plugin"
)

func main() {
	p , _ = plugin.Open("hexify.so")
	f := p.Lookup("Hexigy")
	f.Println(f.(func(string) string)("Nice to meet you"))
	//f.Println(f.("hello"))
}

//go build -buildmode=shared hexify.go
