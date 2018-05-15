// pair (value, type)

package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x float64 = 3.14

	v := reflect.ValueOf(x) // x is a value

	fmt.Println(v)

	vp := reflect.ValueOf(&x) // x is a pointer

	fmt.Println(vp)

	ve := vp.Elem() // value
	fmt.Println(ve)

	fmt.Println(ve.CanSet())
	fmt.Println(ve.CanAddr())
	fmt.Println(ve.CanInterface())

	ve.SetFloat(6.14)
	fmt.Println(ve.Interface())
	fmt.Println(x) // At last x is be modified

}
