// 测试 otto.New()

package main

import (
	"fmt"

	"github.com/robertkrimen/otto"
)

var vm = otto.New()

func main() {
	todo1()
	todo2()
	todo3()
}

func todo1() {
	// vm := otto.New()
	vm.Set("num", 1)
	// s := `console.log("Num: " + num);`
	// s := "console.log(\"Num: \" + num);"
	script, _ := vm.Compile("", "console.log(\"Num: \" + num);")
	v, _ := vm.Run(script)
	// v,_ := vm.Run(s)
	fmt.Println("vvvvv: ", v)
}

func todo2() {
	// vm := otto.New()

	vm.Set("num", 2)
	vm.Run(`
			console.log("Num: " + num);
		`)
}

func todo3() {
	// vm := otto.New()
	vm.Run(`
			console.log("Num: " + num);
		`)
	vm.Set("num", 3)
	vm.Run(`
			console.log("Num: " + num);
		`)
}
