package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

// 嵌套结构体
type ss struct {
	as struct {
		age  int
		name string
	}
	age  int
	name string
	bool
	float64
}

func (s ss) Method1(i int) string  { return "结构体方法1" }
func (s *ss) Method2(i int) string { return "结构体方法2" }

var (
	intValue   = int(0)
	int8Value  = int8(8)
	int16Value = int16(16)
	int32Value = int32(32)
	int64Value = int64(64)

	uIntValue   = uint(0)
	uInt8Value  = uint8(8)
	uInt16Value = uint16(16)
	uInt32Value = uint32(32)
	uInt64Value = uint64(64)

	byteValue    = byte(0)
	runeValue    = rune(0)
	uintptrValue = uintptr(0)

	boolValue   = false
	stringValue = ""

	float32Value = float32(32)
	float64Value = float64(64)

	complex64Value  = complex64(64)
	complex128Value = complex128(128)

	arrayValue  = [5]string{}           // 数组
	sliceValue  = []byte{0, 0, 0, 0, 0} // 切片
	mapValue    = map[string]int{}      // 映射
	chanValue   = make(chan int, 2)     // 通道
	structValue = ss{                   // 结构体
		struct {
			age  int
			name string
		}{10, "子结构体"},
		20,
		"结构体",
		false,
		64.0,
	}

	func1Value = func(a, b, c int) string { // 函数（固定参数）
		return fmt.Sprintf("固定参数：%v %v %v", a, b, c)
	}

	func2Value = func(a, b int, c ...int) string { // 函数（动态参数）
		return fmt.Sprintf("动态参数：%v %v %v", a, b, c)
	}

	unsafePointer     = unsafe.Pointer(&structValue)    // 通用指针
	reflectType       = reflect.TypeOf(0)               // 反射类型
	reflectValue      = reflect.ValueOf(0)              // 反射值
	reflectArrayValue = reflect.ValueOf([]int{1, 2, 3}) // 切片反射值
	// 反射接口类型
	interfaceType = reflect.TypeOf(new(interface{})).Elem()
)

// 简单类型
var simpleTypes = []interface{}{
	intValue, &intValue, // int
	int8Value, &int8Value, // int8
	int16Value, &int16Value, // int16
	int32Value, &int32Value, // int32
	int64Value, &int64Value, // int64
	uIntValue, &uIntValue, // uint
	uInt8Value, &uInt8Value, // uint8
	uInt16Value, &uInt16Value, // uint16
	uInt32Value, &uInt32Value, // uint32
	uInt64Value, &uInt64Value, // uint64
	byteValue, &byteValue, // byte
	runeValue, &runeValue, // rune
	uintptrValue, &uintptrValue, // uintptr
	boolValue, &boolValue, // bool
	stringValue, &stringValue, // string
	float32Value, &float32Value, // float32
	float64Value, &float64Value, // float64
	complex64Value, &complex64Value, // complex64
	complex128Value, &complex128Value, // complex128
}

// 复杂类型
var complexTypes = []interface{}{
	arrayValue, &arrayValue, // 数组
	sliceValue, &sliceValue, // 切片
	mapValue, &mapValue, // 映射
	chanValue, &chanValue, // 通道
	structValue, &structValue, // 结构体
	func1Value, &func1Value, // 定参函数
	func2Value, &func2Value, // 动参函数
	structValue.Method1, structValue.Method2, // 方法
	unsafePointer, &unsafePointer, // 指针
	reflectType, &reflectType, // 反射类型
	reflectValue, &reflectValue, // 反射值
	interfaceType, &interfaceType, // 接口反射类型
}

// 空值
var unsafeP unsafe.Pointer

// 空接口
var nilInterfece interface{}

func main() {
	// 测试简单类型
	for i := 0; i < len(simpleTypes); i++ {
		PrintInfo(simpleTypes[i])
	}
	// 测试复杂类型
	for i := 0; i < len(complexTypes); i++ {
		PrintInfo(complexTypes[i])
	}
	// 测试单个对象
	PrintInfo(unsafeP)
	PrintInfo(&unsafeP)
	PrintInfo(nilInterfece)
	PrintInfo(&nilInterfece)
}

func PrintInfo(i interface{}) {
	if i == nil {
		fmt.Println("--------------------")
		fmt.Printf("无效接口值：%v\n", i)
		return
	}
	t := reflect.TypeOf(i) // 第一步都是通过TypOf(i)获取type

	PrintType(t)
}

func PrintType(t reflect.Type) {
	fmt.Println("--------------------")
	// ----- 通用方法 -----
	fmt.Println("String             :", t.String())     // 类型字符串
	fmt.Println("Name               :", t.Name())       // 类型名称
	fmt.Println("PkgPath            :", t.PkgPath())    // 所在包名称
	fmt.Println("Kind               :", t.Kind())       // 所属分类
	fmt.Println("Size               :", t.Size())       // 内存大小
	fmt.Println("Align              :", t.Align())      // 字节对齐
	fmt.Println("FieldAlign         :", t.FieldAlign()) // 字段对齐
	fmt.Println("NumMethod          :", t.NumMethod())  // 方法数量

	if t.NumMethod() > 0 {
		i := 0
		for ; i < t.NumMethod()-1; i++ {
			fmt.Println("    ┣", t.Method(i).Name) // 通过索引定位方法
		}
		fmt.Println("    ┗", t.Method(i).Name) // 通过索引定位方法
	}
	if sm, ok := t.MethodByName("String"); ok { // 通过名称定位方法
		fmt.Println("MethodByName       :", sm.Index, sm.Name)
	}
	fmt.Println("Implements(i{})    :", t.Implements(interfaceType))  // 是否实现了指定接口
	fmt.Println("ConvertibleTo(int) :", t.ConvertibleTo(reflectType)) // 是否可转换为指定类型
	fmt.Println("AssignableTo(int)  :", t.AssignableTo(reflectType))  // 是否可赋值给指定类型的变量
	fmt.Println("Comparable         :", t.Comparable())               // 是否可进行比较操作
	// ----- 特殊类型 -----
	switch t.Kind() {
	// 指针型：
	case reflect.Ptr:
		fmt.Println("=== 指针型 ===")
		// 获取指针所指对象
		t = t.Elem()
		fmt.Printf("转换到指针所指对象 : %v\n", t.String())
		// 递归处理指针所指对象
		PrintType(t)
		return
	// 自由指针型：
	case reflect.UnsafePointer:
		fmt.Println("=== 自由指针 ===")
		// ...
	// 接口型：
	case reflect.Interface:
		fmt.Println("=== 接口型 ===")
		// ...
	}
	// ----- 简单类型 -----
	// 数值型：
	if reflect.Int <= t.Kind() && t.Kind() <= reflect.Complex128 {
		fmt.Println("=== 数值型 ===")
		fmt.Println("Bits               :", t.Bits()) // 位宽
	}
	// ----- 复杂类型 -----
	switch t.Kind() {
	// 数组型：
	case reflect.Array:
		fmt.Println("=== 数组型 ===")
		fmt.Println("Len                :", t.Len())  // 数组长度
		fmt.Println("Elem               :", t.Elem()) // 数组元素类型
	// 切片型：
	case reflect.Slice:
		fmt.Println("=== 切片型 ===")
		fmt.Println("Elem               :", t.Elem()) // 切片元素类型
	// 映射型：
	case reflect.Map:
		fmt.Println("=== 映射型 ===")
		fmt.Println("Key                :", t.Key())  // 映射键
		fmt.Println("Elem               :", t.Elem()) // 映射值类型
	// 通道型：
	case reflect.Chan:
		fmt.Println("=== 通道型 ===")
		fmt.Println("ChanDir            :", t.ChanDir()) // 通道方向
		fmt.Println("Elem               :", t.Elem())    // 通道元素类型
	// 结构体：
	case reflect.Struct:
		fmt.Println("=== 结构体 ===")
		fmt.Println("NumField           :", t.NumField()) // 字段数量
		if t.NumField() > 0 {
			var i, j int
			// 遍历结构体字段
			for i = 0; i < t.NumField()-1; i++ {
				field := t.Field(i) // 获取结构体字段
				fmt.Printf("    ├ %v\n", field.Name)
				// 遍历嵌套结构体字段
				if field.Type.Kind() == reflect.Struct && field.Type.NumField() > 0 {
					for j = 0; j < field.Type.NumField()-1; j++ {
						subfield := t.FieldByIndex([]int{i, j}) // 获取嵌套结构体字段
						fmt.Printf("    │    ├ %v\n", subfield.Name)
					}
					subfield := t.FieldByIndex([]int{i, j}) // 获取嵌套结构体字段
					fmt.Printf("    │    └ %%v\n", subfield.Name)
				}
			}
			field := t.Field(i) // 获取结构体字段
			fmt.Printf("    └ %v\n", field.Name)
			// 通过名称查找字段
			if field, ok := t.FieldByName("ptr"); ok {
				fmt.Println("FieldByName(ptr)   :", field.Name)
			}
			// 通过函数查找字段
			if field, ok := t.FieldByNameFunc(func(s string) bool { return len(s) > 3 }); ok {
				fmt.Println("FieldByNameFunc    :", field.Name)
			}
		}
	// 函数型：
	case reflect.Func:
		fmt.Println("=== 函数型 ===")
		fmt.Println("IsVariadic         :", t.IsVariadic()) // 是否具有变长参数
		fmt.Println("NumIn              :", t.NumIn())      // 参数数量
		if t.NumIn() > 0 {
			i := 0
			for ; i < t.NumIn()-1; i++ {
				fmt.Println("    ┣", t.In(i)) // 获取参数类型
			}
			fmt.Println("    ┗", t.In(i)) // 获取参数类型
		}
		fmt.Println("NumOut             :", t.NumOut()) // 返回值数量
		if t.NumOut() > 0 {
			i := 0
			for ; i < t.NumOut()-1; i++ {
				fmt.Println("    ┣", t.Out(i)) // 获取返回值类型
			}
			fmt.Println("    ┗", t.Out(i)) // 获取返回值类型
		}
	}
}
