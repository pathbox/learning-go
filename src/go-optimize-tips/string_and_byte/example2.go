//转换后 [ ]byte 底层数组与原 string 内部指针并不相同，以此可确定数据被复制。
// 那么，如不修改数据，仅转换类型，是否可避开复制，从而提升性能？

// 从 ptype 输出的结构来看，string 可看做 [2]uintptr，
// 而 [ ]byte 则是 [3]uintptr，这便于我们编写代码，无需额外定义结构类型。
// 如此，str2bytes 只需构建 [3]uintptr{ptr, len, len}，
// 而 bytes2str 更简单，直接转换指针类型，忽略掉 cap 即可。

package main

import (
	"fmt"
	"strings"
	"unsafe"
)

func string2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func bytes2string(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func main() {
	s := strings.Repeat("abc", 3)
	b := string2bytes(s)
	s2 := bytes2string(b)

	fmt.Println(b, s2)
}

// h
// type = struct []uint8 {
//     uint8 *array;
//     int len;
//     int cap;
// }

// x
// type = struct string {
//     uint8 *str;
//     int len;
// }
