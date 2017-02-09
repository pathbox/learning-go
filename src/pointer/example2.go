package main

import (
	"fmt"
)

type Test struct {
	Name string
}

func change2(t *Test) {
	t.Name = "2"
}

func change3(t *Test) {
	//注意这里括号
	//如果直接*t.Name=3 编译不通过 报错 invalid indirect of t.Name (type string)
	//其实在go里面*可以省掉，直接类似change2函数里这样使用。
	(*t).Name = "3"
}

func change4(t Test) {
	t.Name = "4"
}

func main() {
	// t 是一个地址
	t := &Test{Name: "1"}
	fmt.Println(t.Name)

	change2(t)
	fmt.Println(t.Name)

	change3(t)
	fmt.Println(t.Name)

	// 这里传递变量用了*
	change4(*t)
	fmt.Println(t.Name)
}
