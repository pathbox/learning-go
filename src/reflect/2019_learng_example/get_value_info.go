// 获取各种值的相关信息
package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

// 嵌套结构体
type ss struct {
	a struct {
		int
		string
	}
	int
	string
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
	stringValue = "stringValue"

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
			int
			string
		}{10, "子结构体"},
		20,
		"结构体",
		false,
		64.0,
	}

	func1Value = func(i int) string { // 函数（固定参数）
		return fmt.Sprintf("固定参数：%v", i)
	}

	func2Value = func(i ...int) string { // 函数（动态参数）
		return fmt.Sprintf("动态参数：%v", i)
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
	PrintInfo(&unsafeP)
	PrintInfo(nilInterfece)
	// PrintInfo(&nilInterfece) // 会引发 panic
}

func PrintInfo(i interface{}) {
	if i == nil {
		fmt.Println("--------------------")
		fmt.Printf("无效接口值：%v\n", i)
		fmt.Println("--------------------")
		return
	}
	v := reflect.ValueOf(i)
	PrintValue(v)
}

func PrintValue(v reflect.Value) {
	fmt.Println("--------------------")
	// ----- 通用方法 -----
	fmt.Println("String             :", v.String())  // 反射值的字符串形式
	fmt.Println("Type               :", v.Type())    // 反射值的类型
	fmt.Println("Kind               :", v.Kind())    // 反射值的类别
	fmt.Println("CanAddr            :", v.CanAddr()) // 是否可以获取地址
	fmt.Println("CanSet             :", v.CanSet())  // 是否可以修改
	if v.CanAddr() {
		fmt.Println("Addr               :", v.Addr())       // 获取地址
		fmt.Println("UnsafeAddr         :", v.UnsafeAddr()) // 获取自由地址
	}
	// 是否可转换为接口对象
	fmt.Println("CanInterface       :", v.CanInterface())
	if v.CanInterface() {
		fmt.Println("Interface          :", v.Interface()) // 转换为接口对象
	}
	// 获取方法数量
	fmt.Println("NumMethod          :", v.NumMethod())
	if v.NumMethod() > 0 {
		// 遍历方法
		i := 0
		for ; i < v.NumMethod()-1; i++ {
			fmt.Printf("    ┣ %v\n", v.Method(i).String())
			//			if i >= 4 { // 只列举 5 个
			//				fmt.Println("    ┗ ...")
			//				break
			//			}
		}
		fmt.Printf("    ┗ %v\n", v.Method(i).String())
		// 通过名称获取方法
		fmt.Println("MethodByName       :", v.MethodByName("String").String())
	}
	// ----- 可获取指针的类型 -----
	switch v.Kind() {
	case reflect.Slice, reflect.Map, reflect.Chan, reflect.Func,
		reflect.Ptr, reflect.UnsafePointer:
		fmt.Println("Pointer            :", v.Pointer())
	}
	// ----- 特殊类型 -----
	switch v.Kind() {
	// 指针：
	case reflect.Ptr:
		fmt.Println("=== 指针 ===")
		// 获取指针地址
		if !v.IsNil() {
			// 获取指针所指对象
			v = v.Elem() // 只有指针和接口类型可以使用 Elem()
			fmt.Printf("转换到指针所指对象 : %v\n", v.Type())
			// 递归处理指针所指对象
			PrintValue(v)
			return
		}
	// 自由指针：
	case reflect.UnsafePointer:
		fmt.Println("=== 自由指针 ===")
		if v.Pointer() == 0 {
			v.SetPointer(unsafePointer)
			fmt.Println("重新指向新对象     :", v.Pointer())
		}
		// 将自由指针转换为 *ss 指针（因为定义 unsafePointer 时已经确定了类型）
		s := (*ss)(v.Interface().(unsafe.Pointer))
		// 获取反射值
		v = reflect.ValueOf(s)
		if !v.IsNil() {
			// 获取指针所指对象
			v = v.Elem() // 只有指针和接口类型可以使用 Elem()
			fmt.Printf("转换到指针所指对象 : %v\n", v.Type())
			// 递归处理指针所指对象
			PrintValue(v)
			return
		}
	// 接口：
	case reflect.Interface:
		fmt.Println("=== 接口 ===")
		// 获取接口数据
		fmt.Println("InterfaceData      :", v.InterfaceData())
		// 获取接口所包含的对象
		v = v.Elem() // 只有指针和接口类型可以使用 Elem()
		fmt.Printf("转换到接口所含对象 : %v\n", v.Type())
		// 递归处理接口的动态对象
		PrintValue(v)
		return
	}
	// ----- 简单类型 -----
	// 有符号整型：
	if reflect.Int <= v.Kind() && v.Kind() <= reflect.Int64 {
		fmt.Println("=== 有符号整型 ===")
		fmt.Println("Int                :", v.Int()) // 获取值
		if v.CanSet() {
			v.SetInt(10)                                 // 设置值
			fmt.Println("Int                :", v.Int()) // 获取值
			v.Set(reflect.ValueOf(20).Convert(v.Type())) // 设置值
		}
		fmt.Println("Int                :", v.Int())           // 获取值
		fmt.Println("OverflowInt        :", v.OverflowInt(10)) // 是否溢出
	}
	// 无符号整型：
	if reflect.Uint <= v.Kind() && v.Kind() <= reflect.Uint64 {
		fmt.Println("=== 无符号整型 ===")
		fmt.Println("Uint               :", v.Uint()) // 获取值
		if v.CanSet() {
			v.SetUint(10)                                 // 设置值
			fmt.Println("Uint               :", v.Uint()) // 获取值
			v.Set(reflect.ValueOf(20).Convert(v.Type()))  // 设置值
		}
		fmt.Println("Uint               :", v.Uint())           // 获取值
		fmt.Println("OverflowUint       :", v.OverflowUint(10)) // 是否溢出
	}
	switch v.Kind() {
	// 浮点数：
	case reflect.Float32, reflect.Float64:
		fmt.Println("=== 浮点数 ===")
		fmt.Println("Float              :", v.Float()) // 获取值
		if v.CanSet() {
			v.SetFloat(10)                                 // 设置值
			fmt.Println("Float              :", v.Float()) // 获取值
			v.Set(reflect.ValueOf(20).Convert(v.Type()))   // 设置值
		}
		fmt.Println("Float              :", v.Float())           // 获取值
		fmt.Println("OverflowFloat      :", v.OverflowFloat(10)) // 是否溢出
	// 复数：
	case reflect.Complex64, reflect.Complex128:
		fmt.Println("=== 复数 ===")
		fmt.Println("Complex            :", v.Complex()) // 获取值
		if v.CanSet() {
			v.SetComplex(10)                                   // 设置值
			fmt.Println("Complex            :", v.Complex())   // 获取值
			v.Set(reflect.ValueOf(20 + 20i).Convert(v.Type())) // 设置值
		}
		fmt.Println("Complex            :", v.Complex())           // 获取值
		fmt.Println("OverflowComplex    :", v.OverflowComplex(10)) // 是否溢出
	// 布尔型：
	case reflect.Bool:
		fmt.Println("=== 布尔型 ===")
		fmt.Println("Bool               :", v.Bool()) // 获取值
		if v.CanSet() {
			v.SetBool(true)                               // 设置值
			fmt.Println("Bool               :", v.Bool()) // 获取值
			v.Set(reflect.ValueOf(false))                 // 设置值
		}
		fmt.Println("Bool               :", v.Bool()) // 获取值
	// 字符串：
	case reflect.String:
		fmt.Println("=== 字符串 ===")
		fmt.Println("String             :", v.String()) // 获取值
		if v.CanSet() {
			v.SetString("abc")                              // 设置值
			fmt.Println("String             :", v.String()) // 获取值
			v.Set(reflect.ValueOf("def"))                   // 设置值
		}
		fmt.Println("String             :", v.String()) // 获取值
	// ----- 复杂类型 -----
	// 切片型：
	case reflect.Slice:
		fmt.Println("=== 切片型 ===")
		fmt.Println("Len                :", v.Len()) // 获取长度
		fmt.Println("Cap                :", v.Cap()) // 获取容量
		if v.CanSet() {
			v.SetLen(4) // 不能大于 cap
			v.SetCap(4) // 不能小于 len，只能缩，不能扩
			fmt.Println("SetLen, SetCap     :", v.Len(), v.Cap())
			// 重新指定字节内容
			if v.Type().Elem().Kind() == reflect.Uint8 {
				v.SetBytes([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0})
			}
			fmt.Println("SetByte            :", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0})
		}
		// 获取字节内容
		if v.Type().Elem().Kind() == reflect.Uint8 {
			fmt.Println("Bytes              :", v.Bytes())
		}
		// 根据索引获取元素
		if v.Len() > 0 {
			for i := 0; i < v.Len(); i++ {
				fmt.Println("Index              :", v.Index(i))
			}
		}
		// 获取一个指定范围的切片
		// 参数：起始下标，结束下标
		// 长度 = 结束下标 - 起始下标
		s1 := v.Slice(1, 2)
		fmt.Println("Slice              :", s1)
		fmt.Println("Len                :", s1.Len()) // 获取长度
		fmt.Println("Cap                :", s1.Cap()) // 获取容量
		// 获取一个指定范围和容量的切片
		// 参数：起始下标，结束下标，容量下标
		// 长度 = 结束下标 - 起始下标
		// 容量 = 容量下标 - 起始下标
		s2 := v.Slice3(1, 2, 4)
		fmt.Println("Slice              :", s2)
		fmt.Println("Len                :", s2.Len()) // 获取长度
		fmt.Println("Cap                :", s2.Cap()) // 获取容量
	// 映射型：
	case reflect.Map:
		fmt.Println("=== 映射型 ===")
		// 设置键值，不需要检测 CanSet
		v.SetMapIndex(reflect.ValueOf("a"), reflect.ValueOf(1))
		v.SetMapIndex(reflect.ValueOf("b"), reflect.ValueOf(2))
		v.SetMapIndex(reflect.ValueOf("c"), reflect.ValueOf(3))
		// 获取键列表
		fmt.Println("MapKeys            :", v.MapKeys())
		for _, idx := range v.MapKeys() {
			// 根据键获取值
			fmt.Println("MapIndex           :", v.MapIndex(idx))
		}
	// 结构体：
	case reflect.Struct:
		fmt.Println("=== 结构体 ===")
		// 获取字段个数
		fmt.Println("NumField           :", v.NumField())
		if v.NumField() > 0 {
			var i, j int
			// 遍历结构体字段
			for i = 0; i < v.NumField()-1; i++ {
				field := v.Field(i) // 获取结构体字段
				fmt.Printf("    ├ %-8v %v\n", field.Type(), field.String())
				// 遍历嵌套结构体字段
				if field.Kind() == reflect.Struct && field.NumField() > 0 {
					for j = 0; j < field.NumField()-1; j++ {
						subfield := v.FieldByIndex([]int{i, j}) // 获取嵌套结构体字段
						fmt.Printf("    │    ├ %-8v %v\n", subfield.Type(), subfield.String())
						// if i >= 4 { // 只列举 5 个
						// 	fmt.Println("        ┗ ...")
						// 	break
						// }
					}
					subfield := v.FieldByIndex([]int{i, j}) // 获取嵌套结构体字段
					fmt.Printf("    │    └ %-8v %v\n", subfield.Type(), subfield.String())
				}
				// if i >= 4 { // 只列举 5 个
				// 	fmt.Println("    ┗ ...")
				// 	break
				// }
			}
			field := v.Field(i) // 获取结构体字段
			fmt.Printf("    └ %-8v %v\n", field.Type(), field.String())
			// 通过名称查找字段
			if v := v.FieldByName("ptr"); v.IsValid() {
				fmt.Println("FieldByName(ptr)   :", v.Type().Name())
			}
			// 通过函数查找字段
			v := v.FieldByNameFunc(func(s string) bool { return len(s) > 3 })
			if v.IsValid() {
				fmt.Println("FieldByNameFunc    :", v.Type().Name())
			}
		}
	// 通道型：
	case reflect.Chan:
		fmt.Println("=== 通道型 ===")
		// 发送数据（会阻塞）
		v.Send(reflectValue)
		// 尝试发送数据（不会阻塞）
		fmt.Println("TrySend            :", v.TrySend(reflectValue))
		// 接收数据（会阻塞）
		if x, ok := v.Recv(); ok {
			fmt.Println("Recv               :", x) //
		}
		// 尝试接收数据（不会阻塞）
		if x, ok := v.TryRecv(); ok {
			fmt.Println("TryRecv            :", x) //
		}
		// 因为要执行两次，通道和通道指针各执行一次，关闭后第二次就无法执行了。
		// v.Close()
	// 函数型：
	case reflect.Func:
		fmt.Println("=== 函数型 ===")
		// 判断函数是否具有变长参数
		if v.Type().IsVariadic() {
			// 与可变参数对应的实参必须是切片类型的反射值（reflectArrayValue）。
			fmt.Println("CallSlice          :", v.CallSlice([]reflect.Value{reflectArrayValue})) //
			// 也可以用 v.Call 调用变长参数的函数，只需传入 reflectValue 即可。
		} else {
			// 根据函数定义的参数数量，传入相应数量的反射值（reflectValue）。
			fmt.Println("Call               :", v.Call([]reflect.Value{reflectValue})) //
		}
	}
}
