package main

import "fmt"

type Magician interface {
  Magic()
}

type Delegate struct {
  magician Magician   // struct magician　就是Magicain　接口类型
}

func (delegate Delegate) Magic() {
  if _, ok := delegate.magician.(Magician); ok {
    delegate.magician.Magic()
      fmt.Println("magic")
  }else {
    fmt.Println("base magic")
  }
}

func (delegate Delegate) MoreMagic(){
  delegate.magician.Magic()
	delegate.magician.Magic()
}

type Foo struct {
  Delegate
}

// func (foo Foo) Magic() {
// 	fmt.Println("foo magic")
// }

func main() {
  f := Delegate{new(Foo)}
  f.Magic() // print magic
  f.MoreMagic()
}
