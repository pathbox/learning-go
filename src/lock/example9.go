package main

import (
	"log"
	"sync"
)

// 将匿名函数中 wg 的传入类型改为 *sync.WaitGrou,这样就能引用到正确的WaitGroup了。

// func main() {
// 	wg := &sync.WaitGroup{}
// 	for i := 0; i < 100; i++ {
// 		wg.Add(1)
// 		go func(wg *sync.WaitGroup, i int) {
// 			log.Printf("i: %d", i)
// 			wg.Done()
// 		}(wg, i)
// 	}
// 	wg.Wait()
// 	log.Println("exit")
// }

// 将匿名函数中的 wg 的传入参数去掉，因为Go支持闭包类型，在匿名函数中可以直接使用外面的 wg 变量
func main() {
	wg := sync.WaitGroup{}
	// wg := &sync.WaitGroup{}  // 这是等价的
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			log.Printf("i: %d", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Println("exit")
}
