package main

import (
	"fmt"

	"github.com/willf/bitset"
)

func main() {
	var b bitset.BitSet                  // 定义一个BitSet对象
	b.Set(10).Set(11).Set(100).Set(1000) // 给这个set新增两个值10和11
	if b.Test(1000) {                    // 查看set中是否有1000这个值（我觉得Test这个名字起得是真差劲，为啥不叫Exist）
		b.Clear(1000) // 情况set
	}
	for i, e := b.NextSet(0); e; i, e = b.NextSet(i + 1) { // 遍历整个Set
		fmt.Println("The following bit is set:", i)
	}
	if b.Intersection(bitset.New(100).Set(10)).Count() > 1 { // set求交集
		fmt.Println("Intersection works.")
	}
}
