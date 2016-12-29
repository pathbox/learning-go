package main

import (
	"fmt"
)

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

func zoo() {
	defer func() {
		fmt.Println("zoo defer func invoked")
	}()

	fmt.Println("zoo invoked")
	panic("runtime exception")
}

func main() {
	foo()
}

// 从结果可以看出：
//     panic在zoo中发生，在zoo真正退出前，zoo中注册的defer函数会被逐一执行(FILO)，由于zoo defer中没有捕捉panic，因此panic被抛向其caller：bar。
//     这时对于bar而言，其函数体中的zoo的调用就好像变成了panic调用似的，zoo有些类似于“黑客帝国3”中里奥被史密斯(panic)感 染似的，也变成了史密斯(panic)。panic在bar中扩展开来，bar中的defer也没有捕捉和recover panic，因此在bar中的defer func执行完毕后，panic继续抛给bar的caller: foo；
//     这时对于foo而言，bar就变成了panic，同理，最终foo将panic抛给了main
//     main与上述函数一样，没有recover，直接异常返回，导致进程异常退出。
