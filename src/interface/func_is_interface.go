package main

import "log"

type Job interface { // 定义一个Job接口
	Run() // method
}

type Schedule struct { // schedule struct
	Job  Job // Job 为Job接口
	Name string
}

type JobFunc func() // JobFunc 为 func() 类型

func main() {
	s := &Schedule{
		Name: "Say Hello",
	}

	j := func() { // 定义一个func()，可以理解为先注册了一个拥有具体逻辑的func()
		log.Println("Hello World!")
	}

	jf := JobFunc(j) // 很巧妙的方式！ 这一步将 普通的func 转为了JobFunc() 类型的 func() 这样，jf就是JobFunc类型，实现了Job interface 方法

	s.Job = jf // 赋值给 s.Job

	s.Job.Run() // 利用s.Job.Run() 执行Run方法，这样真正执行的就是 j的代码逻辑

}

func (j JobFunc) Run() { // JobFunc 实现 interface的Run()方法,这样就实现了interface类型
	j()
}

//------------------------------------------------------------
// 下面是 能够传参数的方法，这样定义的各个func 都要有对应的参数个数和类型

// package main

// import "log"

// type Job interface { // 定义一个Job接口
// 	Run(string) // method
// }

// type Schedule struct { // schedule struct
// 	Job  Job // Job 为Job接口
// 	Name string
// }

// type JobFunc func(string) // JobFunc 为 func() 类型

// func main() {
// 	s := &Schedule{
// 		Name: "Say Hello",
// 	}

// 	j := func(s string) { // 定义一个func()
// 		log.Println("Hello World!" + s)
// 	}

// 	jf := JobFunc(j) // 很巧妙的方式！ 这一步将 普通的func 转为了JobFunc() 类型的 func() 这样，jf就是JobFunc类型，实现了Job interface 方法

// 	s.Job = jf // 赋值给 s.Job

// 	s.Job.Run("Cary") // 利用s.Job.Run() 执行Run方法，这样真正执行的就是 j的代码逻辑

// }

// func (j JobFunc) Run(s string) { // JobFunc 实现 interface的Run()方法,这样就实现了interface类型
// 	j(s)
// }
