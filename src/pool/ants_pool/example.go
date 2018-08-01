package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	ants "github.com/panjf2000/ants"
)

var sum int32

func myFunc(i interface{}) error {
	n := i.(int32)
	atomic.AddInt32(&sum, n)
	fmt.Printf("run with %d\n", n)
	return nil
}

func demoFunc() error {
	time.Sleep(10 * time.Millisecond)
	fmt.Println("Hello World!")
	return nil
}

func main() {
	defer ants.Release()

	runTimes := 1000

	var wg sync.WaitGroup
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		ants.Submit(func() error {
			demoFunc()
			wg.Done()
			return nil
		})
	}

	wg.Wait()
	fmt.Printf("running goroutines: %d\n", ants.Running())
	fmt.Printf("finish all tasks.\n")

	// use the pool with a function
	// set 10 the size of goroutine pool and 1 second for expired duration
	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) error {
		myFunc(i)
		wg.Done()
		return nil
	})
	defer p.Release()

	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		p.Serve(int32(i))
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", p.Running())
	fmt.Printf("finish all tasks, result is %d\n", sum)
}
