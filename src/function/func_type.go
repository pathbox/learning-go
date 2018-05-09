package main

import "fmt"

// 定义一个func类型，定义好参数和返回值，具体这个selfDoHandler的逻辑代码是怎样的，在初始化这个类型具体实例的时候定义
type selfDoHandler func(param string) string

func main() {

	var sdh selfDoHandler //定义一个变量sdh 为 selfDoHandler类型

	sdh = func(p string) string { // 这个func 的参数和返回值和 selfDoHandler 类型一致就可以
		s := "Hello" + p // 具体实现selfDoHandler 方法逻辑
		return s
	}

	fmt.Println(sdh(" World"))
}

/*
可以用在哪？ 定义一个map
map[string]selfDoHandler

根据key，可以设置不同的selfDoHandler，但是他们的参数和返回值都和selfDoHandler一致，只是具体实现的时候的代码逻辑可以不一样

这样不同的key相当于进行了不同的逻辑处理。具体的逻辑实现还可以用定义一个func的方式
*/
