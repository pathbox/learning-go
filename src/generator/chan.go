package main

import (
	"fmt"
	"time"
)

func fib(n int) chan int {
  c := make(chan int)
  go func(){
    x,y := 0, 1
    for x < n {
      c <- x
      x, y = y, x+y
      time.Sleep(1 * time.Second)
    }
    close(c)
  }()
  return c
}

func main() {
  for i := range fib(1000){
    fmt.Println(i)
  }
}
