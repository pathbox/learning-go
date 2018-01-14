package main

import "fmt"

// 有三种属性
type options struct {
	a int64
	b string
	c map[int]string
}

func NewOption(opts ...ServerOption) *options {
	r := &options{}
	for _, o := range opts { //每次传入相同的options
		o(r)
	}
	return r
}

type ServerOption func(*options) // type is a function

func WriteA(s int64) ServerOption {
	return func(o *options) {
		o.a = s
	}
}

func WriteB(s string) ServerOption {
	return func(o *options) {
		o.b = s
	}
}

func WriteC(s map[int]string) ServerOption {
	return func(o *options) {
		o.c = s
	}
}

func main() {
	opt1 := WriteA(int64(1)) // 只是准备了参数，返回闭包，还没真正执行代码
	opt2 := WriteB("test")
	opt3 := WriteC(make(map[int]string, 0))

	aa := &options{}
	ob := WriteA(int64(100)) // 准备参数,返回闭包ob
	ob(aa)                   // 向闭包中传参数(*options)，执行闭包中的代码。
	fmt.Println("aa", aa.a)

	WriteB("nice")(aa)
	fmt.Println("aa.nice", aa.b)

	op := NewOption(opt1, opt2, opt3) // 真正执行 闭包中的代码

	fmt.Println(op.a, op.b, op.c)
}

/*
一个函数有参数传入，代码内容是闭包，闭包也有参数传入，这样的函数需要进行两次参数传入
func myf(name string) {
	return func(u *user) {
		u.name = name
	}
}

f := myf("John")
u := &user
f(u)
u.name => "John"

或者直接 myf("John")(u)
*/
