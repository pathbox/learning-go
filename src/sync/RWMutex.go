package main

import (
	"fmt"
	"sync"
)

func main() {
	m := make(map[int]int)

	ml := &sync.RWMutex{}
	go func() {
		var i int
		for {
			ml.RLock()
			a := m[1]
			fmt.Println(i + a)
			i += 1
			ml.RUnlock()
		}
	}()
	go func() {
		var i int
		for {
			ml.RLock()
			a := m[100]
			fmt.Println(i + a)
			i += 1
			ml.RUnlock()
		}
	}()
	go func() {
		for {
			ml.Lock()
			m[2] = 2
			ml.Unlock()
		}
	}()
	select {}
}

// 读写锁读操作和写操作时要同时使用，要不会发生 panic 报错

// fatal error: concurrent map read and map write
