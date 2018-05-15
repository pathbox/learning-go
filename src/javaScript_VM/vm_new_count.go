package main

import (
	"fmt"
	"time"

	"github.com/robertkrimen/otto"
)

func main() {

	// ary := []*otto.Otto{}
	for i := 0; i < 100000; i++ {
		vm := otto.New()
		// ary = append(ary, vm)
		vm.Set("aassddfdjgkjldflgjfdgjdfgljflgjfgjgfg", "dsfjdiofjsdfljdsfklsjdfljsdfjfkdjl")
	}
	// fmt.Println(len(ary))
	fmt.Println("done")
	time.Sleep(1000 * time.Second)
}

// 1w 个 vm  runtime  1.5G 内存
