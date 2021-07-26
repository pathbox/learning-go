package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("going:", i)
			wg.Wait() // 所有goroutine都阻塞在这里，直到wg.Done()执行了才会继续向下走
			fmt.Println("go done:", i)
		}(i)
	}

	time.Sleep(2 * time.Second)
	fmt.Println("wg done")
	wg.Done()
	time.Sleep(3 * time.Second)

}
