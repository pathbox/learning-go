package main

import (
	"fmt"
)

func main() {
	s := "hello world!"
	b := []byte(s)

	fmt.Println(s, b)
}

// go build -gcflags "-N -l" -o test example1.go

// sudo gdb test

// l main.main

// b 11

// r

// info locals

// ptype s

// ptype b

// x/2xg &s

// x/3xg &b

// quit

//转换后 [ ]byte 底层数组与原 string 内部指针并不相同，以此可确定数据被复制。
// 那么，如不修改数据，仅转换类型，是否可避开复制，从而提升性能？

// 从 ptype 输出的结构来看，string 可看做 [2]uintptr，
// 而 [ ]byte 则是 [3]uintptr，这便于我们编写代码，无需额外定义结构类型。
// 如此，str2bytes 只需构建 [3]uintptr{ptr, len, len}，
// 而 bytes2str 更简单，直接转换指针类型，忽略掉 cap 即可。
