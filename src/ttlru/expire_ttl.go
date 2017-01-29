package main

import (
	"fmt"
	"time"

	"github.com/zvelo/ttlru"
)

func main() {
	cache := ttlru.New(10, 3*time.Second)
	fmt.Println(cache.Cap(), cache.Len())

	cache.NoReset = true

	cache.Set("foo", "bar")

	val, ok := cache.Get("foo")
	fmt.Println(val, ok)

	time.Sleep(2 * time.Second)

	val, ok = cache.Get("foo")
	fmt.Println(val, ok)

	time.Sleep(2 * time.Second)

	val, ok = cache.Get("foo")
	fmt.Println(val, ok)

	time.Sleep(4 * time.Second)

	val, ok = cache.Get("foo")
	fmt.Println(val, ok)
}
