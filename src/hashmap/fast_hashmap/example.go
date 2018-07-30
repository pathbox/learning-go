package main

import (
	"fmt"
	"sync/atomic"

	hashmap "github.com/cornelk/hashmap"
)

func main() {
	m := &hashmap.HashMap{}
	m.Set("amount", 123)

	amount, ok := m.Get("amount")
	if ok {
		fmt.Println("amount: ", amount)
	}

	var i int64
	actual, _ := m.GetOrInsert("api/123", &i)
	counter := (actual).(*int64)

	atomic.AddInt64(counter, 1)

	count := atomic.LoadInt64(counter)

	fmt.Println("count:", count)

}
