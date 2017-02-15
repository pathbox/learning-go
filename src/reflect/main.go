package main

import (
	"fmt"
	"reflect"
)

var typeRegistry = make(map[string]reflect.Type)

func registerType(elem interface{}) {
	t := reflect.TypeOf(elem).Elem()
	typeRegistry[t.Name()] = t
}

func newStruct(name string) (interface{}, bool) {
	elem, ok := typeRegistry[name]
	if !ok {
		return nil, false
	}
	return reflect.New(elem).Elem().Interface(), true
}

func init() {
	registerType((*test)(nil))
}

type test struct {
	Name string
	Sex  int
}

func main() {
	structName := "test"

	s, ok := newStruct(structName)
	if !ok {
		return
	}

	t, ok := s.(test)
	if !ok {
		return
	}
	t.Name = "I am test"
	fmt.Println(t, reflect.TypeOf(t))
}
