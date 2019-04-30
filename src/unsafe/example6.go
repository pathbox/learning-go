package main

import (
	"fmt"
	"unsafe"
)

type Num struct {
	i string
	j int64
}

func main() {
	n := Num{i: "EDDYCJY", j: 1}
	nPointer := unsafe.Pointer(&n)

	niPointer := (*string)(unsafe.Pointer(nPointer)) // 直接取出指针转为Pointer 再强制转为字符串类型指针值
	*niPointer = "鸡翅"

	njPointer := (*int64)(unsafe.Pointer(uintptr(nPointer) + unsafe.Offsetof(n.j))) //j 为第二个成员变量。需要进行偏移量计算，才可以对其内存地址进行修改。在进行了偏移运算后，当前地址已经指向第二个成员变量。接着重复转换赋值即可
	*njPointer = 2

	fmt.Printf("n.i: %s, n.j: %d", n.i, n.j)
}

/*
结构体的成员变量在内存存储上是一段连续的内存

结构体的初始地址就是第一个成员变量的内存地址

基于结构体的成员地址去计算偏移量。就能够得出其他成员变量的内存地址

ptr := uintptr(nPointer) 这样是错误的，uintptr 类型是不能存储在临时变量中的。因为从 GC 的角度来看，uintptr 类型的临时变量只是一个无符号整数，并不知道它是一个指针地址
*/
