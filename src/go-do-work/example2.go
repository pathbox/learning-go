package main

import (
	"fmt"
	"github.com/c3mb0/go-do-work"
	"time"
)

type adder struct {
	count int
}

func (a adder) DoWork() {
	a.count++
	fmt.Print(a.count, " ")
	time.Sleep(2 * time.Second)
}

func main() {
	test := adder{count: 0}
	pool := gdw.NewWorkerPool(3)
	defer pool.Close()
	pool.Add(test, 5)
	time.Sleep(1 * time.Second)
	pool.SetPoolSize(1)
	fmt.Printf("\n%d\n", pool.GetQueueDepth())
	pool.Wait()
	fmt.Println()
}
