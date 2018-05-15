package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Name string
	Age  int
}

func main() {
	u := User{"Cary", 27}
	s := reflect.ValueOf(&u).Elem()
	typeOfT := s.Type() //把s.Type()返回的Type对象复制给typeofT，typeofT也是一个反射, typeOfT.Field() fields 名称集合

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i) //迭代s的各个域，注意每个域仍然是反射
		fmt.Printf("%d: %s %s = %v\n", i, typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
	sr, _ := typeOfT.FieldByName("Name")
	fmt.Println(sr)
	s.Field(0).SetString("Joe")
	s.Field(1).SetInt(25)
	fmt.Println("Now u: ", u)
}
