package main

import (
	"fmt"
)

func xrange() chan int {
  var ch = make(chan int)
  go func(){
    for i:= 0; ; i++ {
      ch <-i
    }
  }()

  return ch
}

func main() {
  g := xrange()
  for i := 0; i < 10; i++{
    fmt.Println(<-g)
  }
}

// 闭包的一个作用效果　闭包中的变量不会新建和重新初始化，而是将上一次的执行后的变量带入到下一次的代码执行中
//　所以　闭包往往和循环或迭代在一起使用才有意义
