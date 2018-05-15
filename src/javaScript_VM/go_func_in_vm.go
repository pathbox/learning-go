package main

import (
	"fmt"

	"github.com/robertkrimen/otto"
)

func main() {
	vm := otto.New()

	vm.Set("sayHello", func(call otto.FunctionCall) otto.Value {
		fmt.Printf("Hello %s.\n", call.Argument(0).String()) // 传入第一个参数
		return otto.Value{}
	})

	vm.Set("twoPlus", func(call otto.FunctionCall) otto.Value {
		right, _ := call.Argument(0).ToInteger()
		left, _ := call.Argument(1).ToInteger()
		result, _ := vm.ToValue(left + right)
		return result
	})

	result, _ := vm.Run(`
			sayHello("World");
			sayHello();

			r = twoPlus(1,2);
		`)

	fmt.Println("The result: ", result)
}
