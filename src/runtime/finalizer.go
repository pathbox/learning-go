package main

import (
	"fmt"
	"runtime"
	"time"
)

type Foo struct {
	name string
	num  int
}

func finalizer(f *Foo) {
	fmt.Println("a finalizer has run for ", f.name, f.num)
}

var counter int

func MakeFoo(name string) (a_foo *Foo) {
	a_foo = &Foo{name, counter}
	counter++
	runtime.SetFinalizer(a_foo, finalizer)
	return
}

func Bar() {
	f1 := MakeFoo("one")
	f2 := MakeFoo("two")

	fmt.Println("f1 is: ", f1.name)
	fmt.Println("f2 is: ", f2.name)
}

func main() {
	for i := 0; i < 3; i++ {
		Bar()
		time.Sleep(time.Second)
		runtime.GC()
	}
	fmt.Println("done.")
}

// 如果需要在一个对象 obj 被从内存移除前执行一些特殊操作，比如写到日志文件中，可以通过如下方式调用函数来实现
