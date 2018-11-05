package main

import (
	"fmt"
	"runtime"
)

func main() {
	Foo()
}
func Foo() {
	fmt.Printf("我是 %s, 谁在调用我?\n", printMyName())
	Bar()
}
func Bar() {
	fmt.Printf("我是 %s, 谁又在调用我?\n", printMyName())
}

func printMyName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
