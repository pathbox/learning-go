package main

import (
	"fmt"
	"reflect"
)

// // 获取 t 类型的字符串描述，不要通过 String 来判断两种类型是否一致。
// func (t *rtype) String() string

// // 获取 t 类型在其包中定义的名称，未命名类型则返回空字符串。
// func (t *rtype) Name() string

// // 获取 t 类型所在包的名称，未命名类型则返回空字符串。
// func (t *rtype) PkgPath() string

// // 获取 t 类型的类别。
// func (t *rtype) Kind() reflect.Kind

// // 获取 t 类型的值在分配内存时的大小，功能和 unsafe.SizeOf 一样。
// func (t *rtype) Size() uintptr

// // 获取 t 类型的值在分配内存时的字节对齐值。
// func (t *rtype) Align() int

// // 获取 t 类型的值作为结构体字段时的字节对齐值。
// func (t *rtype) FieldAlign() int

// // 获取 t 类型的方法数量。
// func (t *rtype) NumMethod() int

// // 根据索引获取 t 类型的方法，如果方法不存在，则 panic。
// // 如果 t 是一个实际的类型，则返回值的 Type 和 Func 字段会列出接收者。
// // 如果 t 只是一个接口，则返回值的 Type 不列出接收者，Func 为空值。
// func (t *rtype) Method() reflect.Method

// // 根据名称获取 t 类型的方法。
// func (t *rtype) MethodByName(string) (reflect.Method, bool)

// // 判断 t 类型是否实现了 u 接口。
// func (t *rtype) Implements(u reflect.Type) bool

// // 判断 t 类型的值可否转换为 u 类型。
// func (t *rtype) ConvertibleTo(u reflect.Type) bool

// // 判断 t 类型的值可否赋值给 u 类型。
// func (t *rtype) AssignableTo(u reflect.Type) bool

// // 判断 t 类型的值可否进行比较操作
// func (t *rtype) Comparable() bool

// 示例
type inf interface {
	Method1()
	Method2()
}

type ss struct {
	a func()
}

func (i ss) Method1() {}
func (i ss) Method2() {}

func main() {
	s := reflect.TypeOf(ss{})
	i := reflect.TypeOf(new(inf)).Elem()

	Test(s)
	Test(i)
}

func Test(t reflect.Type) {
	if t.NumMethod() > 0 {
		fmt.Printf("\n--- %s ---\n", t)
		fmt.Println(t.Method(0).Type)
		fmt.Println(t.Method(0).Func.String())
	}
}
