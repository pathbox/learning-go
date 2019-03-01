package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	num  = 1000 * 10
	gnum = 1000
)

func main() {
	fmt.Println("only read")
	testRwmutexReadOnly()
	testMutexReadOnly()

	fmt.Println("write and read")
	testRwmutexWriteRead()
	testMutexWriteRead()

	fmt.Println("write only")
	testRwmutexWriteOnly()
	testMutexWriteOnly()
}

func testRwmutexReadOnly() {
	var w = &sync.WaitGroup{}
	var rwmutexTmp = newRwmutex()
	w.Add(gnum)
	t1 := time.Now()
	for i := 0; i < gnum; i++ {
		go func() {
			defer w.Done()
			for in := 0; in < num; in++ {
				rwmutexTmp.get(in)
			}
		}()
	}
	w.Wait()
	fmt.Println("testRwmutexReadOnly cost:", time.Now().Sub(t1).String())
}

func testRwmutexWriteOnly() {
	var w = &sync.WaitGroup{}
	var rwmutexTmp = newRwmutex()
	w.Add(gnum)
	t1 := time.Now()
	for i := 0; i < gnum; i++ {
		go func() {
			defer w.Done()
			for in := 0; in < num; in++ {
				rwmutexTmp.set(in, in)
			}
		}()
	}
	w.Wait()
	fmt.Println("testRwmutexWriteOnly cost:", time.Now().Sub(t1).String())
}

func testRwmutexWriteRead() {
	var w = &sync.WaitGroup{}
	var rwmutexTmp = newRwmutex()
	w.Add(gnum)
	t1 := time.Now()
	for i := 0; i < gnum; i++ {
		if i%2 == 0 {
			go func() {
				defer w.Done()
				for in := 0; in < num; in++ {
					rwmutexTmp.get(in)
				}
			}()
		} else {
			go func() {
				defer w.Done()
				for in := 0; in < num; in++ {
					rwmutexTmp.set(in, in)
				}
			}()
		}
	}
	w.Wait()
	fmt.Println("testRwmutexWriteRead cost:", time.Now().Sub(t1).String())
}

func testMutexReadOnly() {
	var w = &sync.WaitGroup{}
	var mutexTmp = newMutex()
	w.Add(gnum)

	t1 := time.Now()
	for i := 0; i < gnum; i++ {
		go func() {
			defer w.Done()
			for in := 0; in < num; in++ {
				mutexTmp.get(in)
			}
		}()
	}
	w.Wait()
	fmt.Println("testMutexReadOnly cost:", time.Now().Sub(t1).String())
}

func testMutexWriteOnly() {
	var w = &sync.WaitGroup{}
	var mutexTmp = newMutex()
	w.Add(gnum)

	t1 := time.Now()
	for i := 0; i < gnum; i++ {
		go func() {
			defer w.Done()
			for in := 0; in < num; in++ {
				mutexTmp.set(in, in)
			}
		}()
	}
	w.Wait()
	fmt.Println("testMutexWriteOnly cost:", time.Now().Sub(t1).String())
}

func testMutexWriteRead() {
	var w = &sync.WaitGroup{}
	var mutexTmp = newMutex()
	w.Add(gnum)
	t1 := time.Now()
	for i := 0; i < gnum; i++ {
		if i%2 == 0 {
			go func() {
				defer w.Done()
				for in := 0; in < num; in++ {
					mutexTmp.get(in)
				}
			}()
		} else {
			go func() {
				defer w.Done()
				for in := 0; in < num; in++ {
					mutexTmp.set(in, in)
				}
			}()
		}

	}
	w.Wait()
	fmt.Println("testMutexWriteRead cost:", time.Now().Sub(t1).String())
}

func newRwmutex() *rwmutex {
	var t = &rwmutex{}
	t.mu = &sync.RWMutex{}
	t.ipmap = make(map[int]int, 100)

	for i := 0; i < 100; i++ {
		t.ipmap[i] = 0
	}
	return t
}

type rwmutex struct {
	mu    *sync.RWMutex
	ipmap map[int]int
}

func (t *rwmutex) get(i int) int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.ipmap[i]
}

func (t *rwmutex) set(k, v int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	k = k % 100
	t.ipmap[k] = v
}

func newMutex() *mutex {
	var t = &mutex{}
	t.mu = &sync.Mutex{}
	t.ipmap = make(map[int]int, 100)

	for i := 0; i < 100; i++ {
		t.ipmap[i] = 0
	}
	return t
}

type mutex struct {
	mu    *sync.Mutex
	ipmap map[int]int
}

func (t *mutex) get(i int) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.ipmap[i]
}

func (t *mutex) set(k, v int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	k = k % 100
	t.ipmap[k] = v
}

// 在读多写少下，sync.Map的性能要比sync.RwMutex + map高的多
