package main

import (
	"fmt"
	"unsafe"
)

func main() {
	s := struct {
		a byte
		b byte
		c byte
		d int64
	}{0, 0, 0, 0}

	// 将结构体指针转换为通用指针
	p := unsafe.Pointer(&s)

	// 保存结构体的地址备用（偏移量为0）
	up0 := uintptr(p) // 结构体s的初始地址

	// 将通用指针转换为 byte 型指针
	pb := (*byte)(p)

	// 给转换后的指针赋值
	*pb = 100

	// 结构体内容跟着改变
	fmt.Println(s)

	// 偏移到第2个字段
	up := up0 + unsafe.Offsetof(s.b)

	// 将偏移后的地址转换为通用指针
	p = unsafe.Pointer(up)

	// 将通用指针转换为 byte 型指针
	pb = (*byte)(p)
	// 给转换后的指针赋值
	*pb = 200
	// 结构体内容跟着改变
	fmt.Println(s)

	// 偏移到第 3 个字段
	up = up0 + unsafe.Offsetof(s.c)

	// 将偏移后的地址转换为通用指针
	p = unsafe.Pointer(up)

	// 将通用指针转换为 byte 型指针
	pb = (*byte)(p)
	// 给转换后的指针赋值  byte 类型是有上限的，255上限
	*pb = 255
	// 结构体内容跟着改变
	fmt.Println(s)

	// 偏移到第 4 个字段
	up = up0 + unsafe.Offsetof(s.d)
	// 将偏移后的地址转换为通用指针
	p = unsafe.Pointer(up)
	// 将通用指针转换为 int64 型指针
	pi := (*int64)(p)
	// 给转换后的指针赋值
	*pi = 40
	// 结构体内容跟着改变
	fmt.Println(s)
}
