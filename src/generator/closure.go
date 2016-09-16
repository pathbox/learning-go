package main

import (
	"fmt"
	// "time"
)


func fib() func() int{
  a, b := 0, 1
  return func() int {
    a, b = b, a+b
    return a
  }
}

func main() {
  f := fib()
  // Function calls are evaluated left-to-right.
	// fmt.Println(f(), f(), f(), f(), f())
  for i := 0; i<14; i++{
    fmt.Println(f())  // fmt.Println(fib()()) 这种方式闭包失效
  }
}

// 这里会执行１４次ｆｉｂ函数。如果没有闭包，得到的结果应该是１
//　由于闭包，a, b := 0, 1只会在第一次的时候执行，第二次的时候是将闭包的变量带入第二次的执行，
//　也就是执行return ｆｕｎｃ()中的代码，并且a,b是上一次的结果
