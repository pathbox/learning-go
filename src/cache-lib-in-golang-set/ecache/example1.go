package main

import (
	"strconv"
	"time"

	"github.com/orca-zhang/ecache"
)

func main() {
	var c = ecache.NewLRUCache(16, 200, 10*time.Second)
	// 整型键
	c.Put(strconv.FormatInt(d, 10), o) // d为`int64`类型

	// 整型值
	c.PutInt64("uid1", int64(1))
	if d, ok := c.GetInt64("uid1"); ok {
		// d为`int64`类型的1
	}

	// 字节数组
	c.PutBytes("uid1", b) // b为`[]byte`类型
	if b, ok := c.GetBytes("uid1"); ok {
		// b为`[]byte`类型
	}
	var c = ecache.NewLRUCache(16, 200, 10*time.Second).LRU2(1024)

	cache.Inspect(func(action int, key string, iface *interface{}, bytes []byte, status int) {
		// TODO: 实现你想做的事情
		//     监听器会根据注入顺序依次执行
		//     注意⚠️如果有耗时操作，尽量另开channel保证不阻塞当前协程

		// - 如何获取正确的值 -
		//   - `Put`:      `*iface`
		//   - `PutBytes`: `bytes`
		//   - `PutInt64`: `ecache.ToInt64(bytes)`
	})
}
