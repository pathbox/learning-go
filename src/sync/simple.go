package main

import (
	"log"
	"sync"
)

func main() {
	// 建立对象
	var pool = &sync.Pool{New: nil}
	// 准备放入的字符串
	val := "Hello,World!"
	// 放入
	pool.Put(val)
	// 取出
	log.Println(pool.Get())
	// 再取就没有了,会自动调用NEW
	log.Println(pool.Get())

	str := "Hello Kitty"
	pool.Put(str)
	log.Println(pool.Get())
	log.Println(pool.Get())
}
