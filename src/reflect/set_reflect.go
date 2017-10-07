// 通过判断反射返回类型，元素是否能被修改，是否存在等成功修改相应等字段值

package main

import (
	"fmt"
	"reflect"
)

type SUser struct {
	Id   int
	Name string
	Age  int
}

func SetReflectInit() {
	u := SUser{1, "Mike", 11}
	SetInfo(&u)
	fmt.Println(u)
}

func SetInfo(o interface{}) {
	v := reflect.ValueOf(o) // 反射获取传入的接口值 这时候获取的是指针 &u

	if v.Kind() == reflect.Ptr && !v.Elem().CanSet() { // 判断是否为指针类型 元素是否可以修改
		fmt.Println("can not set")
		return
	}
	v = v.Elem() // 实际取得的对象  SUser{1, "Mike", 11}

	// 判断字段是否存在 Name 字段是否存在
	f := v.FieldByName("Name")
	if !f.IsValid() {
		fmt.Println("var invalid")
		return
	}

	// 设置字段
	if f := v.FieldByName("Name"); f.Kind() == reflect.String {
		f.SetString("John") // 将原来的 Mike 设置为 John
	}

}

func main() {
	SetReflectInit()
}
