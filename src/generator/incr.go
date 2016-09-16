package main

import "fmt"

func incr(max int) func() int {
  init := 0 // 由于闭包，这行代码只会在第一次的时候执行
  return func() int{
    // if init == max{
    //   init = 0
    // }
    init++
    return init
  }
}
//　闭包作用
func main() {
  f := incr(5)　// 这一句，就是闭包的定义一定要将闭包函数赋值给一个变量
  for i := 0; i < 15; i++  {
    fmt.Println(f())　// 1-15
  }
}

// 闭包失效
func main() {
	for i := 0; i < 15; i++ {
		fmt.Println(incr(5)()) //  1 1 1 1 1...1 1
	}
}
