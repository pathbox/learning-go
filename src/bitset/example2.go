package main

import (
	"fmt"

	"github.com/bits-and-blooms/bitset"
)

func main() {
	b := bitset.New(64) // 构造一个64bit长度的bitset
	// 放入一个值，第11位是1
	b.Set(10)
	fmt.Println(b.DumpAsBits()) // 64位长度的01字符串
	b.Clear(10)                 // 删除一个值
	b.Set(1).Set(3)
	fmt.Println(b.Len())
	fmt.Println(b.DumpAsBits())
	// 测试
	fmt.Println(b.Test(3)) // true
	fmt.Println(b.Test(4)) // false

	// 指定位置操作
	b = &bitset.BitSet{}
	b.Set(3)
	// // 在指定位置插入0
	b.InsertAt(3)
	fmt.Println(b.DumpAsBits())
	// 在指定位置修改
	b.SetTo(4, false)
	fmt.Println(b.DumpAsBits()) // 000000000000000000000000000000000000000000000000000000000000
	// 指定位置删除
	b.Set(3).DeleteAt(3) // 0000
	// 签到签出 指定位置插入和删除
}
