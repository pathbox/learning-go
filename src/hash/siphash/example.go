package main

import "fmt"
import "github.com/dchest/siphash"

func main() {
	key := []byte("1")
	h := siphash.New(key)
	h.Write([]byte("Hello"))
	// sum := h.Sum(nil)
	sum := h.Sum64()
	fmt.Println(sum)
}
