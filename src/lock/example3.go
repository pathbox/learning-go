// 锁成功的例子
package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	Value int
}

var wg sync.WaitGroup

func main() {
	mutex := &sync.Mutex{} // 定义了一个锁 mutex
	wg.Add(1000)
	counter := &Counter{Value: 0}

	for i := 0; i < 1000; i++ {
		go Count(counter, mutex)
	}
	wg.Wait()
	fmt.Println("Count Value: ", counter.Value)
}

func Count(counter *Counter, mutex *sync.Mutex) {
	mutex.Lock()
	defer mutex.Unlock()
	counter.Value++
	wg.Done()
}

// 结果输出是：Count Value:  1000。 加锁成功
