package main

import (
	"fmt"

	"github.com/bits-and-blooms/bitset"
)

// 两个bitsets交互
func main() {
	a := &bitset.BitSet{}
	a.Set(1).Set(3).Set(5) // 索引位置是从0开始计算的
	fmt.Println(a.DumpAsBits())
	b := &bitset.BitSet{}
	b.Set(3).Set(5).Set(7)
	// 交集
	fmt.Println(a.Intersection(b)) // {3,5}
	// 并集
	fmt.Println(a.Union(b)) // {1,3,5,7}
	// 差集
	fmt.Println(a.Difference(b)) // {1}
	// 全等
	fmt.Println(a.Equal(b)) // false

	b = bitset.New(64)
	b.Set(1).Set(3).Set(5).Set(7)
	for i, e := b.NextSet(0); e; i, e = b.NextSet(i + 1) {
		fmt.Println("The following bit is set:", i, e)
	}
}
