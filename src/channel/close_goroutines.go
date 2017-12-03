package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

const (
	N = 10
)

func main() {
	runtime.GOMAXPROCS(2)
	quit := make(chan bool)

	for i := 0; i < N; i++ {
		go func(name string) {
			for {
				select {
				case <-quit:
					fmt.Printf("clean up %s\n", name)
					return
				}
			}
		}(strconv.Itoa(i))
	}
	close(quit) // 会通知所有goroutine，quit <- true 这只会通知其中一个goroutine

	time.Sleep(3 * time.Second)
}
