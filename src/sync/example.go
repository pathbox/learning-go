package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

func main() {
	// 禁用GC,保证在main函数执行结束前恢复GC
	defer debug.SetGCPercent(debug.SetGCPercent(-1))
	var count int32
	newFunc := func() interface{} {
		return atomic.AddInt32(&count, 1)
	}
	pool := sync.Pool{New: newFunc}

	// New 字段值的作用
	v1 := pool.Get()
	fmt.Printf("v1: %v\n", v1)

	// 临时对象池的存取
	pool.Put(newFunc())
	pool.Put(newFunc())
	pool.Put(newFunc())
	v2 := pool.Get()
	fmt.Printf("v2: %v\n", v2)

	// 垃圾回收对临时对象池的影响
	debug.SetGCPercent(100)
	runtime.GC()
	v3 := pool.Get()
	fmt.Printf("v3: %v\n", v3)
	// pool.New = nil
	v4 := pool.Get()
	fmt.Printf("v4: %v\n", v4)
}

/*

在实现过程中还要特别注意的是Pool本身也是一个对象，要把Pool对象在程序开始的时候初始化为全局唯一。
对象池使用是较简单的，但原生的sync.Pool有个较大的问题：我们不能自由控制Pool中元素的数量，放进Pool中的对象每次GC发生时都会被清理掉。这使得sync.Pool做简单的对象池还可以，但做连接池就有点心有余而力不足了，比如：在高并发的情景下一旦Pool中的连接被GC清理掉，那每次连接DB都需要重新三次握手建立连接，这个代价就较大了。

*/
