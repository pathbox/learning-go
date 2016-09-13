package main

import (
	"fmt"
	"time"

	"github.com/bluele/gcache"
)

func main() {
	gc := gcache.New(20).
		LRU().
		Expiration(2 * time.Second).
		LoaderFunc(func(key interface{}) (interface{}, error) {
			return "ok1", nil
		}).
		Build()
	gc.Set("key", "ok")
	value, err := gc.Get("key")
	if err != nil {
		panic(err)
	}
	fmt.Println("Get:", value)

	time.Sleep(3 * time.Second)

	value1, err := gc.Get("key")
	if err != nil {
		panic(err)
	}
	fmt.Println("Get:", value1)
}
