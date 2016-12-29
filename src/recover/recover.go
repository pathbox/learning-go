package main

import (
	"fmt"
)

func zoo() {
	defer func() {
		fmt.Println("zoo defer func1 invoked")
	}()

	defer func() {
		if x := recover(); x != nil {
			fmt.Printf("recover panic: %v in zoo recover defer func", x)
		}
	}()

	defer func() {
		fmt.Println("zoo defer func2 invoked")
	}()

	fmt.Println("zoo invoked")
	panic("zoo runtime exception")
}

func foo() {
	defer func() {
		fmt.Println("foo defer func invoked")
	}()
	fmt.Println("foo invoked")

	bar()
	fmt.Println("do something after bar in foo")
}

func bar() {
	defer func() {
		fmt.Println("bar defer func invoked")
	}()
	fmt.Println("bar invoked")

	zoo()
	fmt.Println("do something after zoo in bar")
}

func main() {
	foo()
}

// recover只有在defer函数中调用才能起到recover的作用，这样recover就和defer函数有了紧密联系。我们在zoo的defer函数中捕捉并recover这个panic：

// 由于zoo在defer里恢复了panic，这样在zoo返回后，bar不会感知到任何异常，将按正常逻辑输出函数执行内容，比如：“do something after zoo in bar”,以此类推。
