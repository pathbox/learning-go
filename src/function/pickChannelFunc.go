package main

import (
	"fmt"
	"math/rand"
	"time"
)

var delay = 200 * time.Millisecond

func pickFunc(fns ...func()) func() {
	return fns[rand.Intn(len(fns))]
}

func produce(c chan func(), n int, fns ...func()) {
	defer close(c)
	for i := 0; i < n; i++ {
		c <- pickFunc(fns...)
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	x := 10
	fns := []func(){
		func() { x += 1 },
		func() { x += 1 },
		func() { x *= 2 },
		func() { x /= 2 },
		func() { x *= x },
	}

	c := make(chan func())
	go produce(c, 10, fns...)
	for fn := range c {
		fn()
		fmt.Println(x)
		time.Sleep(delay)
	}
}
