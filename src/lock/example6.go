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
var mutex sync.Mutex // 声明了一个全局锁
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
	mutex.Lock()
	defer mutex.Unlock()
	counter.Value++
	wg.Done()
}

// 结果输出是：Count Value:  1000。 加锁成功
// 总结： 声明一个全局锁，别将其当成参数传入方法中使用
