package main

import (
	"fmt"
	"sync"
)

func main() {
	syncMap := sync.Map{}
	var keys []string

	syncMap.Store("a", "1")
	syncMap.Store("b", "1")
	syncMap.Store("c", "1")
	syncMap.Store("d", "1")
	syncMap.Store("e", "1")

	f := func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	}

	syncMap.Range(f)

	fmt.Println(keys)
}
