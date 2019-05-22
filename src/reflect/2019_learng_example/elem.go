package main

import (
	"fmt"
	"reflect"
)

type Foo struct {
	A int `tag1:"Tag1" tag2:"Second Tag"`
	B string
}

func main() {
	// Struct
	f := Foo{A: 10, B: "Salutations"}
	// Struct类型的指针
	fPtr := &f
	// Map
	m := map[string]int{"A": 1, "B": 2}
	// channel
	ch := make(chan int)
	// slice
	sl := []int{1, 32, 34}
	//string
	str := "string var"
	// string 指针
	strPtr := &str

	tMap := examiner(reflect.TypeOf(f), 0)
	tMapPtr := examiner(reflect.TypeOf(fPtr), 0)
	tMapM := examiner(reflect.TypeOf(m), 0)
	tMapCh := examiner(reflect.TypeOf(ch), 0)
	tMapSl := examiner(reflect.TypeOf(sl), 0)
	tMapStr := examiner(reflect.TypeOf(str), 0)
	tMapStrPtr := examiner(reflect.TypeOf(strPtr), 0)

	fmt.Println("tMap :", tMap)
	fmt.Println("tMapPtr: ", tMapPtr)
	fmt.Println("tMapM: ", tMapM)
	fmt.Println("tMapCh: ", tMapCh)
	fmt.Println("tMapSl: ", tMapSl)
	fmt.Println("tMapStr: ", tMapStr)
	fmt.Println("tMapStrPtr: ", tMapStrPtr)
}

// 类型以及元素的类型判断
func examiner(t reflect.Type, depth int) map[int]map[string]string {
	outType := make(map[int]map[string]string)

	// 如果是一下类型，重新验证
	switch t.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		fmt.Println("这几种类型Name是空字符串：", t.Name(), ", Kind是：", t.Kind())
		// 递归查询元素类型
		tMap := examiner(t.Elem(), depth)
		for k, v := range tMap {
			outType[k] = v
		}

	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i) // reflect字段
			outType[i] = map[string]string{
				"Name": f.Name,
				"Kind": f.Type.String(),
			}
		}
	default:
		// 直接验证类型
		outType = map[int]map[string]string{depth: {"Name": t.Name(), "Kind": t.Kind().String()}}
	}

	return outType
}

/*
Kind()
			Kind有slice, map , pointer指针，struct, interface, string , Array, Function, int或其他基本类型组成。Kind和Type之前要做好区分。如果你定义一个 type Foo struct {}， 那么Kind就是struct,  Type就是Foo。

Elem()
				如果你的变量是一个指针、map、slice、channel、Array。那么你可以使用reflect.Typeof(v).Elem()来确定包含的类型。

*/
