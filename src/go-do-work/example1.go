package main

import (
	"fmt"
	"github.com/c3mb0/go-do-work"
)

type adder struct {
	count int
}

func (a *adder) DoWork() {
	a.count++
	fmt.Print(a.count, " ")
}

func main() {
	test := adder{count: 0}
	pool := gdw.NewWorkerPool(2)
	defer pool.Close()
	pool.Add(&test, 5)
	pool.Wait()
	fmt.Println()
}
