package main

import (
	"fmt"
	"hash/fnv"
)

func main() {
	h := fnv.New64a()
	h.Write([]byte("Hello World"))
	fmt.Println(h.Sum64())
}
