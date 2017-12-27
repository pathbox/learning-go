// 对于小对象，直接将数据交由 map 保存，远比用指针高效。这不但减少了堆内存分配，
// 关键还在于垃圾回收器不会扫描非指针类型 key/value 对象

package main

import (
	"runtime"
	"time"
)

const capacity = 500000

var d interface{}

func value() interface{} {
	m := make(map[int]int, capacity)

	for i := 0; i < capacity; i++ {
		m[i] = i
	}

	return m
}

func pointer() interface{} {
	m := make(map[int]*int, capacity)

	for i := 0; i < capacity; i++ {
		v := i
		m[i] = &v
	}

	return m
}

func main() {
	d = pointer() // d = value()

	for i := 0; i < 20; i++ {
		runtime.GC()
		time.Sleep(time.Second)
	}

}

// map 对 key、value 数据存储长度有限制 128长度限制
