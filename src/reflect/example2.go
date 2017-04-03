package main

import (
	"fmt"
	"reflect"
)

type Home struct {
	i int `nljb:"100"`
}

func main() {
	home := new(Home)
	home.i = 5
	rcvr := reflect.ValueOf(home)
	fmt.Println(rcvr)
	typ := reflect.Indirect(rcvr).Type()
	fmt.Println(typ)
	fmt.Println(typ.Kind().String())
	x := typ.NumField()
	for i := 0; i < x; i++ {
		nljb := typ.Field(0).Tag.Get("nljb") // tag nljb:"100"
		fmt.Println(nljb)
	}
}
