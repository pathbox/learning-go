// map 不会收缩 “不再使用” 的空间。就算把所有键值删除，
// 它依然保留内存空间以待后用,直到后续的GC有必要将这内存回收

package main

import (
	"runtime/debug"
	"time"
)

const capacity = 1000000

var dict = make(map[int][100]byte, capacity)

func test() {
	for i := 0; i < capacity; i++ {
		dict[i] = [100]byte{}
	}

	for k := range dict {
		delete(dict, k)
	}

	dict = nil // 释放map对象,这样才能回收内存
}

func main() {
	test()

	for i := 0; i < 20; i++ {
		debug.FreeOSMemory()
		time.Sleep(time.Second)
	}
}

// 如长期使用 map 对象（比如用作 cache 容器），偶尔换成 “新的” 或许会更好。还有，int key 要比 string key 更快。
// go build -o test && GODEBUG="gctrace=1" ./test
