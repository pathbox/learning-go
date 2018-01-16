package main

import(
	"sync"
	"fmt"
)

// 闭包思想
// worker is a func()
type worker struct {
  Func func()
}

func main() {
	var wg sync.WaitGroup

	channels := make(chan worker, 10) // a simple pool

	 for i := 0; i < 5; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for ch := range channels {
                //reflect.ValueOf(ch.Func).Call(ch.Args)
                ch.Func() // 执行闭包代码
            }
        }()
		}

		for i := 0; i < 100; i++{
			j := i
			wk := worker{
				Func: func() {
					fmt.Println(i+j)
				},
			}
			channels <- wk
		}
		close(channels)
		wg.Wait()

}