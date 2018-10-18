package main

import (
	"fmt"
	"reflect"
	"runtime"
)

func main() {
	fn := say_hello
	name := getFunctionName(fn)
	fmt.Println("function name: ", name)
}

func say_hello() {
	fmt.Println("Hello World")
}

func getFunctionName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf((fn)).Pointer()).Name()
}
