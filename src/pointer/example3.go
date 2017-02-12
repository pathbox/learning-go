package main

import (
	"fmt"
)

type abc struct {
	v int
}

func (a abc) aaaa() { // 传入的是值，而不是引用
	a.v = 1
	fmt.Printf("1:%d\n", a.v)
}

func (a *abc) bbbb() { //传入的是引用，而不是值
	fmt.Printf("2:%d\n", a.v)
	a.v = 2
	fmt.Printf("3:%d\n", a.v)
}

func (a *abc) cccc() { //传入的是引用，而不是值
	a.v = 3
	fmt.Printf("4:%d\n", a.v)
}

func main() {
	aobj := abc{} // new(abc)  aobj 是一个值
	// aobj := new(abc)
	fmt.Println("aobj:", aobj)

	aobj.aaaa()
	fmt.Println("aobj:", aobj)

	aobj.bbbb()
	fmt.Println("aobj:", aobj)

	aobj.cccc()
	fmt.Println("aobj:", aobj)

}
