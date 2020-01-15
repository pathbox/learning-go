package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	bs := []byte{71, 111, 111, 100, 32, 109, 111, 114, 110, 105, 110, 103} // []byte("Good morning")
	s := BytesToString(bs)
	fmt.Println(s)

	b := StringToBytes(s)
	fmt.Println(b)
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func StringToBytes(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}

func NewPtrVal(defValue interface{}) interface{} {
	p := reflect.New(reflect.TypeOf(defValue))
	p.Elem().Set(reflect.ValueOf(defValue))
	return p.Interface()
}
