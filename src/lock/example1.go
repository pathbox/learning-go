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

	counter := &Counter{Value: 0}
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go Count(counter)
	}
	wg.Wait()
	fmt.Println("Count Value: ", counter.Value)
}

func Count(counter *Counter) {
	var mutex sync.Mutex

	mutex.Lock()
	defer mutex.Unlock()
	counter.Value++
	wg.Done()
}

// 结果输出是：Count Value:  982。 982 这个数值是随机的。也就是 Count Value: '小于等于1000的整数'。
// 可以发现，在count 方法中的mutex锁，并没有起到想要的作用。
