package main

import (
	"fmt"
	"time"

	"github.com/koding/cache"
)

func main() {
  cache := cache.NewMemoryWithTTL(3 * time.Second)
  cache.StartGC(time.Second * 3)
  time.Sleep(3 * time.Second)
  cache.Set("test_key", "test_data")

  data, err := cache.Get("test_key")
  fmt.Println(data, err)

  time.Sleep(2 * time.Second)

  data, err = cache.Get("test_key")
  fmt.Println(data, err)
  time.Sleep(2 * time.Second)

  data, err = cache.Get("test_key")
  fmt.Println(data, err)
}
