package main

import (
	"fmt"
)

func main() {
	var p *int = nil
	var i interface{} = p

	fmt.Printf("i: %v \n", i)
	fmt.Printf("i: %x\n", i)
	if i != nil {
		fmt.Println("i is not nil")
	} else {
		fmt.Println("i is nil")
	}
}

/* result:
i: <nil>
i is not nil

i 到底是nil还是不是呢？ 这是 go 中interface的一个坑
虽然我们把一个nil值赋值给interface{}，但是实际上interface里依然存了指向类型的指针，所以拿这个interface变量去和nil常量进行比较的话就会返回false。
这个interface变量其实不是真的nil. 想要避开这个Go语言的坑，我们要做的就是避免将一个有可能为nil的具体类型的值赋值给interface变量。
url: http://studygolang.com/articles/10635#reply3
*/
