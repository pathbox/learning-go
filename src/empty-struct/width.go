package main

import (
	"fmt"
	"unsafe"
)

func main() {
  var s string
  var c complex128
  fmt.Println(unsafe.Sizeof(s))
  fmt.Println(unsafe.Sizeof(c))

  var a [3]uint32
  fmt.Println(unsafe.Sizeof(a))

  type SS struct {
    a uint16
    b uint32
  }
  var ss SS
  fmt.Println(unsafe.Sizeof(ss))
  var sss struct{}
  fmt.Println(unsafe.Sizeof(sss))

  type SSSS struct {
    A struct{}
    B struct{}
  }
  var ssss SSSS
  fmt.Println(unsafe.Sizeof(ssss))
  var x [10000000]struct{}
  fmt.Println(unsafe.Sizeof(x))

  var xx = make([]struct{}, 100000000)
  fmt.Println(unsafe.Sizeof(xx))
}
