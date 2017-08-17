// 锁成的例子
package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	Value int
	sync.Mutex
}

var wg sync.WaitGroup

func main() {

	wg.Add(1000)
	counter := &Counter{Value: 0}

	for i := 0; i < 1000; i++ {
		go Count(counter)
	}
	wg.Wait()
	fmt.Println("Count Value: ", counter.Value)
}

func Count(counter *Counter) {
	counter.Lock()
	defer counter.Unlock()
	counter.Value++
	wg.Done()
}

// 结果输出是：Count Value:  1000。 加锁成功
// 在 Counter struct 上定义了一个属性，这个属性的类型就是互斥锁 sync.Mutex

// 定义的方式还可以有多种

/* 1.
type Counter struct {
   Value int
   Mutex sync.Mutex
}

counter := &Counter{Value: 0}
counter.Mutex.Lock()
defer counter.Mutex.Unlock()

2.
type Counter struct {
   Value int
   Mutex *sync.Mutex
}

counter := &Counter{Value: 0, Mutex: &sync.Mutex{}}
counter.Mutex.Lock()
defer counter.Mutex.Unlock()
*/
