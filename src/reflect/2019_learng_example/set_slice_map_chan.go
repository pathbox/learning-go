package main

import (
	"fmt"
	"reflect"
)

func main() {
	intSlice := make([]int, 0)
	mapStringInt := make(map[string]int)
	sliceType := reflect.TypeOf(intSlice)
	mapType := reflect.TypeOf(mapStringInt)

	// 创建新值
	intSliceReflect := reflect.MakeSlice(sliceType, 0, 0)
	mapReflect := reflect.MakeMap(mapType)

	// 使用新创建的变量
	v := 10
	rv := reflect.ValueOf(v)
	intSliceReflect = reflect.Append(intSliceReflect, rv) //Append
	intSlice2 := intSliceReflect.Interface().([]int)
	fmt.Println(intSlice2)

	k := "hello"
	rk := reflect.ValueOf(k)
	mapReflect.SetMapIndex(rk, rv) //SetMapIndex
	mapStringInt2 := mapReflect.Interface().(map[string]int)
	fmt.Println(mapStringInt2)
}
