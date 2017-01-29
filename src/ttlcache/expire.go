package main

import (
	"fmt"
	"time"

	"github.com/wunderlist/ttlcache"
)

func main() {
	cache := ttlcache.NewCache(2 * time.Second)
	cache.Set("key", "value")
	value3, exists3 := cache.Get("key")
	fmt.Println(value3, exists3)

	time.Sleep(3 * time.Second)
	value, exists := cache.Get("key")
	fmt.Println(value, exists)

	time.Sleep(3 * time.Second)
	value1, exists1 := cache.Get("key")
	fmt.Println(value1, exists1)

	time.Sleep(3 * time.Second)
	value2, exists2 := cache.Get("key")
	fmt.Println(value2, exists2)

	count := cache.Count()
	fmt.Println(count)
}
