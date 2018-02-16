package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type TestObj struct {
	field1 string
}

type emptyInterface struct {
	typ  *struct{}
	word unsafe.Pointer
}

func main() {
	struct_ := &TestObj{}
	structInter := (interface{})(struct_)
	structPtr := (*emptyInterface)(unsafe.Pointer(&structInter)).word
	field, _ := reflect.TypeOf(structInter).Elem().FieldByName("field1")
	field1Ptr := uintptr(structPtr) + field.Offset
	*((*string)(unsafe.Pointer(field1Ptr))) = "hello"
	fmt.Println(struct_)
}
