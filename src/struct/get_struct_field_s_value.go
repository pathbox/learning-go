package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string
	Age  int
	City string
}

func main() {
	p := Person{
		Name: "Joe",
		Age:  27,
		City: "Beijing",
	}

	value := reflect.ValueOf(p)
	for i := 0; i < value.NumField(); i++ {

		fmt.Printf("Field %d: %v\n", i, value.Field(i))
	}
}
