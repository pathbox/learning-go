package main

import (
	"fmt"
	"unsafe"
)

type People interface {
	Show()
}

type Student struct{}

func (stu *Student) Show() {

}

func live() People {
	var stu *Student
	return stu
}

func main() {
	if live() == nil {
		fmt.Println("AAAAAAA")
	} else {
		fmt.Println("BBBBBBB")
	}
}

// BBBBBBB

//  如果改成　var live interface{} 结果是 AAAAAAA

//可以把 inferface 类型的变量 理解为 一个“容器”，该“容器“”将包含赋值给该变量的值的
// 实际 类型信息 和 值信息，只有inferface变量的 类型信息 和 值信息 两部分都为nil，该接口才被称为 nil接口

// 上面例子中，　live()的值为nil,但是 类型信息中 有方法集,这不为nil,所以 总的 live()不认为是nil

// interface 底层结构

type eface struct { //空接口
	_type *_type         //类型信息
	data  unsafe.Pointer //指向数据的指针(go语言中特殊的指针类型unsafe.Pointer类似于c语言中的void*)
}
type iface struct { //带有方法的接口
	tab  *itab          //存储type信息还有结构实现方法的集合
	data unsafe.Pointer //指向数据的指针(go语言中特殊的指针类型unsafe.Pointer类似于c语言中的void*)
}
type _type struct {
	size       uintptr //类型大小
	ptrdata    uintptr //前缀持有所有指针的内存大小
	hash       uint32  //数据hash值
	tflag      tflag
	align      uint8    //对齐
	fieldalign uint8    //嵌入结构体时的对齐
	kind       uint8    //kind 有些枚举值kind等于0是无效的
	alg        *typeAlg //函数指针数组，类型实现的所有方法
	gcdata     *byte
	str        nameOff
	ptrToThis  typeOff
}
type itab struct {
	inter  *interfacetype //接口类型
	_type  *_type         //结构类型
	link   *itab
	bad    int32
	inhash int32
	fun    [1]uintptr //可变大小 方法集合
}
