package main

import (
	"errors"
	"fmt"
	"reflect"
)

// contain 容器, target 值是否在容器中
func Contain(target interface{}, contain interface{}) (bool, error) {
	conValue := reflect.ValueOf(contain)
	switch reflect.TypeOf(contain).Kind() { // 容器的类型
	case reflect.Slice, reflect.Array: // Array or Slice 类型
		for i := 0; i < conValue.Len(); i++ {
			if conValue.Index(i).Interface() == target {
				return true, nil
			}
		}
	case reflect.Map:
		if conValue.MapIndex(reflect.ValueOf(target)).IsValid() {
			return true, nil
		}
	}

	return false, errors.New("not in contain")
}

func main() {
	testMap()

	testArray()
	testSlice()
}

func testArray() {
	a := 1
	b := [3]int{1, 2, 3}

	fmt.Println(Contain(a, b))

	c := "a"
	d := [4]string{"b", "c", "d", "a"}
	fmt.Println(Contain(c, d))

	e := 1.1
	f := [4]float64{1.2, 1.3, 1.1, 1.4}
	fmt.Println(Contain(e, f))

	g := 1
	h := [4]interface{}{2, 4, 6, 1}
	fmt.Println(Contain(g, h))

	i := [4]int64{}
	fmt.Println(Contain(a, i))
}

func testSlice() {
	a := 1
	b := []int{1, 2, 3}

	fmt.Println(Contain(a, b))

	c := "a"
	d := []string{"b", "c", "d", "a"}
	fmt.Println(Contain(c, d))

	e := 1.1
	f := []float64{1.2, 1.3, 1.1, 1.4}
	fmt.Println(Contain(e, f))

	g := 1
	h := []interface{}{2, 4, 6, 1}
	fmt.Println(Contain(g, h))

	i := []int64{}
	fmt.Println(Contain(a, i))
}

func testMap() {
	var a = map[int]string{1: "1", 2: "2"}
	fmt.Println(Contain(3, a))

	var b = map[string]int{"1": 1, "2": 2}
	fmt.Println(Contain("1", b))

	var c = map[string][]int{"1": {1, 2}, "2": {2, 3}}
	fmt.Println(Contain("6", c))
}
