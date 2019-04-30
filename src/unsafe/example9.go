package main

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

func main() {
	// 创建一个 strings 包中的 Reader 对象
	// 它有三个私有字段：s string、i int64、prevRune int
	sr := strings.NewReader("abcdef")
	// 此时 sr 中的成员是无法修改的
	fmt.Println(sr)
	b1, err := sr.ReadByte() // 从i 开始read出
	fmt.Printf("%c, %v\n", b1, err)
	// 但是我们可以通过 unsafe 来进行修改
	// 先将其转换为通用指针
	p := unsafe.Pointer(sr)
	// 获取结构体地址
	up0 := uintptr(p)
	// 确定要修改的字段（这里不能用 unsafe.Offsetof 获取偏移量，因为是私有字段）
	if sf, ok := reflect.TypeOf(*sr).FieldByName("i"); ok { // 修改i的值
		// 偏移到指定字段的地址
		up := up0 + sf.Offset
		// 转换为通用指针
		p = unsafe.Pointer(up)
		// 转换为相应类型的指针
		pi := (*int64)(p)
		// 对指针所指向的内容进行修改
		*pi = 3 // 修改索引
	}
	// 看看修改结果
	fmt.Println(sr)
	// 看看读出的是什么
	b, err := sr.ReadByte() // 从i 开始read出
	fmt.Printf("%c, %v\n", b, err)
}
