package main

import (
	"fmt"
	"time"

	"github.com/bluele/gcache"
)
func main() {
  gc := gcache.New(2).
    EvictedFunc(func(key, value interface{}){
      fmt.Println("evicted key:", key)
    }).
    Build()
    for i := 0; i < 3; i++ {
		gc.Set(i, i*i)
	}
	time.Sleep(2 * time.Second)
}
