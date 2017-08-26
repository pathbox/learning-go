package main

import (
	"fmt"
	"sync"
	"time"
)

type Lane struct {
	mu sync.Mutex
	m  map[byte]byte
}

func testMutexLane(Go, Gnt int) {
	var lanes [256]Lane
	for i := range lanes {
		lanes[i].m = make(map[byte]byte)
	}
	var wg sync.WaitGroup
	wg.Add(Go)

	start := time.Now()
	for g := 0; g < Go; g++ {
		go func() {
			for i := 0; i < Gnt; i++ {
				index := byte(i)
				l := &lanes[index]
				l.mu.Lock()
				l.m[index] = index
				_ = l.m[index]
				l.mu.Unlock()
			}
			wg.Done()
		}()
	}

	wg.Wait()
	end := time.Now()
	use := end.Sub(start)
	op := use / time.Duration(Go*Gnt)
	fmt.Printf("Times=%10v, Go=%4v, Gnt=%8v, Use=%12v %10v/op\n",
		Go*Gnt, Go, Gnt, use, op)
}

func main() {
	testMutexLane(1, 10000000)
	testMutexLane(2, 10000000)
	testMutexLane(4, 10000000)
	testMutexLane(8, 10000000)
	testMutexLane(16, 1000000)
	testMutexLane(32, 1000000)
	testMutexLane(64, 1000000)
	testMutexLane(128, 100000)
	testMutexLane(512, 100000)
	testMutexLane(1014, 10000)
	testMutexLane(2048, 10000)
}
