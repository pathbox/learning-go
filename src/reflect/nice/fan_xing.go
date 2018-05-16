package main

import (
	"fmt"
	"reflect"
	"strings"
)

// 处理不同类型这里是 int 或 string 的add函数, 参数和返回值都是reflect.Value类型
func add(args []reflect.Value) (results []reflect.Value) {
	if len(args) == 0 {
		return nil
	}

	var r reflect.Value
	switch args[0].Kind() {
	case reflect.Int:
		n := 0
		for _, a := range args {
			n += int(a.Int())
		}
		r = reflect.ValueOf(n)
	case reflect.String:
		ss := make([]string, 0, len(args))
		for _, s := range args {
			ss = append(ss, s.String())
		}
		r = reflect.ValueOf(strings.Join(ss, ""))
	}
	results = append(results, r)
	return
}

// wrap reflect.MakeFunc(), T 函数指针会指向add函数体，这样就获得了add的函数体，T相当于传入参数和得到返回值，真正的代码逻辑由add定义
func makeAdd(T interface{}) {
	fn := reflect.ValueOf(T).Elem()
	v := reflect.MakeFunc(fn.Type(), add) // 把原始函数变量的类型和通用算法函数存到同一个value中
	fn.Set(v)                             // 把原始函数指针变量指向v，这样它就获得了函数体
}

func main() {
	// 定义函数变量，未定义函数体
	var intAdd func(x, y, z int) int //相当于传入参数和得到返回值，真正的代码逻辑由add定义
	var strAdd func(a, b string) string

	makeAdd(&intAdd)
	makeAdd(&strAdd)

	fmt.Println(intAdd(12, 23, 33))
	fmt.Println(strAdd("Hello", " World"))
}
