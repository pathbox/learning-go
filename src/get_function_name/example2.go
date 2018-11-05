package main

import (
	"fmt"
	"runtime"
)

func main() {
	Foo()
}
func Foo() {
	fmt.Printf("我是 %s, %s 在调用我!\n", printMyName(), printCallerName())
	Bar()
}
func Bar() {
	fmt.Printf("我是 %s, %s 又在调用我!\n", printMyName(), printCallerName())
}
func printMyName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
func printCallerName() string { // 打印这个函数在哪里被调用，可以用于日志打印报错的时候，就知道是哪个函数报错调用打印日志了
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name()
}
