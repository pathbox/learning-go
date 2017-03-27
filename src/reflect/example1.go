package main

import (
	"fmt"
	"reflect"
)

type MyStruct struct {
	name string
}

func (this *MyStruct) GetName(str string) string {
	this.name = str
	return this.name
}

func main() {
	// 备注: reflect.Indirect -> 如果是指针则返回 Elem()
	// 首先，reflect包有两个数据类型我们必须知道，一个是Type，一个是Value。
	// Type就是定义的类型的一个数据类型，Value是值的类型

	// 对象
	s := "this is string"
	// 获取对象类型 (string)
	fmt.Println("1 ", reflect.TypeOf(s))
	// 获取对象值 (this is string)
	fmt.Println("2 ", reflect.ValueOf(s))
	// 对象

	var x float64 = 3.4

	fmt.Println(reflect.ValueOf(x))
	a := &MyStruct{name: "Jerry"}

	// 返回对象的方法的数量
	fmt.Println("3 ", reflect.TypeOf(a).NumMethod())

	// 遍历对象中的方法
	for m := 0; m < reflect.TypeOf(a).NumMethod(); m++ {
		method := reflect.TypeOf(a).Method(m)
		fmt.Println(method.Type)         // func(*main.MyStruct) string
		fmt.Println(method.Name)         // GetName
		fmt.Println(method.Type.NumIn()) // 参数个数
		fmt.Println(method.Type.In(1))   // 参数类型
	}

	// 获取对象值 (<*main.MyStruct Value>)
	fmt.Println("4 ", reflect.ValueOf(a))
	// 获取对象名称
	fmt.Println("5 ", reflect.Indirect(reflect.ValueOf(a)).Type().Name())

	// 参数
	i := "Hello"
	v := make([]reflect.Value, 0)
	v = append(v, reflect.ValueOf(i))
	// 通过对象值中的方法名称调用方法 ([Jerry]) (返回数组因为Go支持多值返回)
	fmt.Println("6 ", reflect.ValueOf(a).MethodByName("GetName").Call(v))

	// 通过对值中的子对象名称获取值 (Jerry)
	fmt.Println("7 ", reflect.Indirect(reflect.ValueOf(a)).FieldByName("name"))
	// 是否可以改变这个值 (false)
	fmt.Println("8 ", reflect.Indirect(reflect.ValueOf(a)).FieldByName("name").CanSet())
	// 是否可以改变这个值 (true)
	fmt.Println("9 ", reflect.Indirect(reflect.ValueOf(&(a.name))).CanSet())
	// 不可改变 (false)
	fmt.Println("10 ", reflect.Indirect(reflect.ValueOf(s)).CanSet())
	// 可以改变
	// reflect.Indirect(reflect.ValueOf(&s)).SetString("jbnl")
	fmt.Println("11 ", reflect.Indirect(reflect.ValueOf(&s)).CanSet())
}

// 1  string
// 2  this is string
// 3.4
// 3  1
// func(*main.MyStruct, string) string
// GetName
// 2
// string
// 4  &{Jerry}
// 5  MyStruct
// 6  [Hello]
// 7  Hello
// 8  false
// 9  true
// 10  false
// 11  true
