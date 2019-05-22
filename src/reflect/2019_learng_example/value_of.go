package main

import (
	"fmt"
	"reflect"
)

type ss struct {
	A int
	a int
}

func (s ss) Method1(i int) string  { return "结构体方法1" }
func (s *ss) Method2(i int) string { return "结构体方法2" }

func main() {
	v1 := reflect.ValueOf(ss{})                   // 结构体
	v2 := reflect.ValueOf(&ss{})                  // 结构体指针
	v3 := reflect.ValueOf(&ss{}).Elem()           // 可寻址结构体
	v4 := reflect.ValueOf(&ss{}).Elem().Field(0)  // 可寻址结构体的共有字段
	v5 := reflect.ValueOf(&ss{}).Elem().Field(1)  // 可寻址结构体的私有字段
	v6 := reflect.ValueOf(&ss{}).Method(0)        // 结构体指针的方法
	v7 := reflect.ValueOf(&ss{}).Elem().Method(0) // 结构体的方法

	fmt.Println(v1.CanAddr()) // false
	fmt.Println(v2.CanAddr()) // false
	fmt.Println(v3.CanAddr()) // true
	fmt.Println(v4.CanAddr()) // true
	fmt.Println(v5.CanAddr()) // true
	fmt.Println(v6.CanAddr()) // false
	fmt.Println(v7.CanAddr()) // false
	fmt.Println("----------")
	fmt.Println(v1.CanSet()) // false
	fmt.Println(v2.CanSet()) // false
	fmt.Println(v3.CanSet()) // true
	fmt.Println(v4.CanSet()) // true
	fmt.Println(v5.CanSet()) // false
	fmt.Println(v6.CanSet()) // false
	fmt.Println(v6.CanSet()) // false
	fmt.Println("----------")
	fmt.Println(v1.CanInterface()) // true
	fmt.Println(v2.CanInterface()) // true
	fmt.Println(v3.CanInterface()) // true
	fmt.Println(v4.CanInterface()) // true
	fmt.Println(v5.CanInterface()) // false
	fmt.Println(v6.CanInterface()) // true
	fmt.Println(v7.CanInterface()) // true
}

/*
// 将 v 值转换为 uintptr 类型，v 值必须是切片、映射、通道、函数、指针、自由指针。
func (v Value) Pointer() uintptr

// 获取 v 值的地址。v 值必须是可寻址类型（CanAddr）。
func (v Value) UnsafeAddr() uintptr

// 将 UnsafePointer 类别的 v 值修改为 x，v 值必须是 UnsafePointer 类别，必须可修改。
func (v Value) SetPointer(x unsafe.Pointer)

// 判断 v 值是否为 nil，v 值必须是切片、映射、通道、函数、接口、指针。
// IsNil 并不总等价于 Go 的潜在比较规则，比如对于 var i interface{}，i == nil 将返回
// true，但是 reflect.ValueOf(i).IsNil() 将 panic。
func (v Value) IsNil() bool

// 获取“指针所指的对象”或“接口所包含的对象”
func (v Value) Elem() reflect.Value

------------------------------

// 接口

// 获取“指针所指的对象”或“接口所包含的对象”
func (v Value) Elem() reflect.Value

------------------------------

// 通用

// 获取 v 值的字符串描述
func (v Value) String() string

// 获取 v 值的类型
func (v Value) Type() reflect.Type

// 返回 v 值的类别，如果 v 是空值，则返回 reflect.Invalid。
func (v Value) Kind() reflect.Kind

// 获取 v 的方法数量
func (v Value) NumMethod() int

// 根据索引获取 v 值的方法，方法必须存在，否则 panic
// 使用 Call 调用方法的时候不用传入接收者，Go 会自动把 v 作为接收者传入。
func (v Value) Method(int) reflect.Value

// 根据名称获取 v 值的方法，如果该方法不存在，则返回空值（reflect.Invalid）。
func (v Value) MethodByName(string) reflect.Value

// 判断 v 本身（不是 v 值）是否为零值。
// 如果 v 本身是零值，则除了 String 之外的其它所有方法都会 panic。
func (v Value) IsValid() bool

// 将 v 值转换为 t 类型，v 值必须可转换为 t 类型，否则 panic。
func (v Value) Convert(t Type) reflect.Value

------------------------------

// 示例
func main() {
	var v reflect.Value      // 未包含任何数据
	fmt.Println(v.IsValid()) // false

	var i *int
	v = reflect.ValueOf(i)   // 包含一个指针
	fmt.Println(v.IsValid()) // true

	v = reflect.ValueOf(nil) // 包含一个 nil 指针
	fmt.Println(v.IsValid()) // false

	v = reflect.ValueOf(0)   // 包含一个 int 数据
	fmt.Println(v.IsValid()) // true
}

------------------------------

// 获取

// 获取 v 值的内容，如果 v 值不是有符号整型，则 panic。
func (v Value) Int() int64

// 获取 v 值的内容，如果 v 值不是无符号整型（包括 uintptr），则 panic。
func (v Value) Uint() uint64

// 获取 v 值的内容，如果 v 值不是浮点型，则 panic。
func (v Value) Float() float64

// 获取 v 值的内容，如果 v 值不是复数型，则 panic。
func (v Value) Complex() complex128

// 获取 v 值的内容，如果 v 值不是布尔型，则 panic。
func (v Value) Bool() bool

// 获取 v 值的长度，v 值必须是字符串、数组、切片、映射、通道。
func (v Value) Len() int

// 获取 v 值的容量，v 值必须是数值、切片、通道。
func (v Value) Cap() int

// 获取 v 值的第 i 个元素，v 值必须是字符串、数组、切片，i 不能超出范围。
func (v Value) Index(i int) reflect.Value

// 获取 v 值的内容，如果 v 值不是字节切片，则 panic。
func (v Value) Bytes() []byte

// 获取 v 值的切片，切片长度 = j - i，切片容量 = v.Cap() - i。
// v 必须是字符串、数值、切片，如果是数组则必须可寻址。i 不能超出范围。
func (v Value) Slice(i, j int) reflect.Value

// 获取 v 值的切片，切片长度 = j - i，切片容量 = k - i。
// i、j、k 不能超出 v 的容量。i <= j <= k。
// v 必须是字符串、数值、切片，如果是数组则必须可寻址。i 不能超出范围。
func (v Value) Slice3(i, j, k int) reflect.Value

// 根据 key 键获取 v 值的内容，v 值必须是映射。
// 如果指定的元素不存在，或 v 值是未初始化的映射，则返回零值（reflect.ValueOf(nil)）
func (v Value) MapIndex(key Value) reflect.Value

// 获取 v 值的所有键的无序列表，v 值必须是映射。
// 如果 v 值是未初始化的映射，则返回空列表。
func (v Value) MapKeys() []reflect.Value

// 判断 x 是否超出 v 值的取值范围，v 值必须是有符号整型。
func (v Value) OverflowInt(x int64) bool

// 判断 x 是否超出 v 值的取值范围，v 值必须是无符号整型。
func (v Value) OverflowUint(x uint64) bool

// 判断 x 是否超出 v 值的取值范围，v 值必须是浮点型。
func (v Value) OverflowFloat(x float64) bool

// 判断 x 是否超出 v 值的取值范围，v 值必须是复数型。
func (v Value) OverflowComplex(x complex128) bool

------------------------------

// 设置（这些方法要求 v 值必须可修改）

// 设置 v 值的内容，v 值必须是有符号整型。
func (v Value) SetInt(x int64)

// 设置 v 值的内容，v 值必须是无符号整型。
func (v Value) SetUint(x uint64)

// 设置 v 值的内容，v 值必须是浮点型。
func (v Value) SetFloat(x float64)

// 设置 v 值的内容，v 值必须是复数型。
func (v Value) SetComplex(x complex128)

// 设置 v 值的内容，v 值必须是布尔型。
func (v Value) SetBool(x bool)

// 设置 v 值的内容，v 值必须是字符串。
func (v Value) SetString(x string)

// 设置 v 值的长度，v 值必须是切片，n 不能超出范围，不能为负数。
func (v Value) SetLen(n int)

// 设置 v 值的内容，v 值必须是切片，n 不能超出范围，不能小于 Len。
func (v Value) SetCap(n int)

// 设置 v 值的内容，v 值必须是字节切片。x 可以超出 v 值容量。
func (v Value) SetBytes(x []byte)

// 设置 v 值的键和值，如果键存在，则修改其值，如果键不存在，则添加键和值。
// 如果将 val 设置为零值（reflect.ValueOf(nil)），则删除该键。
// 如果 v 值是一个未初始化的 map，则 panic。
func (v Value) SetMapIndex(key, val reflect.Value)

// 设置 v 值的内容，v 值必须可修改，x 必须可以赋值给 v 值。
func (v Value) Set(x reflect.Value)

------------------------------

// 结构体

// 获取 v 值的字段数量，v 值必须是结构体。
func (v Value) NumField() int

// 根据索引获取 v 值的字段，v 值必须是结构体。如果字段不存在则 panic。
func (v Value) Field(i int) reflect.Value

// 根据索引链获取 v 值的嵌套字段，v 值必须是结构体。
func (v Value) FieldByIndex(index []int) reflect.Value

// 根据名称获取 v 值的字段，v 值必须是结构体。
// 如果指定的字段不存在，则返回零值（reflect.ValueOf(nil)）
func (v Value) FieldByName(string) reflect.Value

// 根据匹配函数 match 获取 v 值的字段，v 值必须是结构体。
// 如果没有匹配的字段，则返回零值（reflect.ValueOf(nil)）
func (v Value) FieldByNameFunc(match func(string) bool) Value
*/
