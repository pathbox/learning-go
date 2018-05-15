package main

import (
	"strconv"
	"time"

	"github.com/robertkrimen/otto"
)

func main() {

	vm := otto.New()
	for i := 0; i < 1000000; i++ {
		k := strconv.Itoa(i)
		k = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + k
		vm.Set(k, k)
	}
	time.Sleep(100 * time.Second)
}

// 100w 个 key(size 33+)-value(size 33+)   380M 内存
