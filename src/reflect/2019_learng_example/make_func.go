package main

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

func main() {
	timed := MakeTimedFunction(time1).(func())
	timed()
	timedToo := MakeTimedFunction(time2).(func(int) int)
	time2Val := timedToo(5)
	fmt.Println(time2Val)
}

/*
将创建Func封装， 非reflect.Func类型会panic
当然makeFunc的闭包函数表达式类型是固定的，可以查阅一下文档。
细读文档的reflect.Value.Call()方法。
*/
func MakeTimedFunction(f interface{}) interface{} {
	rf := reflect.TypeOf(f)
	if rf.Kind() != reflect.Func {
		panic("not reflect func")
	}
	vf := reflect.ValueOf(f)
	wrapperF := reflect.MakeFunc(rf, func(in []reflect.Value) []reflect.Value {
		start := time.Now()
		out := vf.Call(in)
		end := time.Now()
		fmt.Printf("calling %s took %v\n", runtime.FuncForPC(vf.Pointer()).Name(), end.Sub(start))
		return out
	})
	return wrapperF.Interface()
}

func time1() {
	fmt.Println("time1Func===starting")
	time.Sleep(1 * time.Second)
	fmt.Println("time1Func===ending")
}

func time2(a int) int {
	fmt.Println("time2Func===starting")
	time.Sleep(time.Duration(a) * time.Second)
	result := a * 2
	fmt.Println("time2Func===ending")
	return result
}
