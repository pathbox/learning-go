package main

import (
	"fmt"
	"reflect"
)

func main() {
	var num float64 = 1.2345

	pointer := reflect.ValueOf(&num)
	value := reflect.ValueOf(num)

	// 可以理解为“强制转换”，但是需要注意的时候，转换的时候，如果转换的类型不完全符合，则直接panic
	// Golang 对类型要求非常严格，类型一定要完全符合
	// 如下两个，一个是*float64，一个是float64，如果弄混，则会panic
	convertPointer := pointer.Interface().(*float64)
	convertValue := value.Interface().(float64)

	fmt.Println(convertPointer)
	fmt.Println(convertValue)

	cType := reflect.TypeOf(convertValue)
	fmt.Println("convertValue Type: ", cType)
	pType := reflect.TypeOf(convertPointer)
	fmt.Println("convertPointer Type: ", pType)

	cValue := reflect.ValueOf(convertValue)
	fmt.Println("convertValue Value: ", cValue)
	pValue := reflect.ValueOf(convertPointer)
	fmt.Println("convertPointer Value: ", pValue)
}
