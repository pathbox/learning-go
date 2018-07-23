package main

import "fmt"

func main() {
	s := "hello world"
	r := hash(s)
	fmt.Println(r)
}

// Jenkins hash function
func hash(key string) uint32 {
	var h uint32
	for _, c := range key {
		h += uint32(c)
		h += (h << 10)
		h ^= (h >> 6)
	}
	h += (h << 3)
	h ^= (h >> 11)
	h += (h << 15)
	return h
}
