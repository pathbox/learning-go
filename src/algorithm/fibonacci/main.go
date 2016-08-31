package main

import "fmt"

func fib(n int) int {
  x, y := 0, 1
  for i := 0; i < n; i++{
    x, y = y, x+y
  }
  return y
}

func main(){
  fmt.Println(fib(10))
}
