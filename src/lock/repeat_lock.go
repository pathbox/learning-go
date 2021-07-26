package main

import (
	"fmt"
	"sync"
	"time"
)

var mutex sync.Mutex

var m map[int]int

func main() {
	m = make(map[int]int)

	for i := 0; i < 100; i++ {
		go func() {
			mutex.Lock()
			m[i] = 1
			mutex.Unlock()
			mutex.Lock()
			m[i] = i
			mutex.Unlock()

		}()
	}

	time.Sleep(5 * time.Second)
	fmt.Println(len(m))
	fmt.Println(m)
}

// 在一个goroutine中mutex锁可以重复使用
