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
		go Count(counter, mutex)
	}
	wg.Wait()
	fmt.Println("Count Value: ", counter.Value)
}

func Count(counter *Counter, mutex sync.Mutex) {
	mutex.Lock()
	defer mutex.Unlock()
	counter.Value++
	wg.Done()
}

// 将mutex在 Count方法外定义，之后再作为参数传入。此时mutex是全局锁
// 结果输出是：Count Value:  982。 982 这个数值是随机的。也就是 Count Value: '小于等于1000的整数'。
// 总结：声明了一个全局锁，作为参数传入到方法中，并不能起到上层范围的互斥锁的作用。
