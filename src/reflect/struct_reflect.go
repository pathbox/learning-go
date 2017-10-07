// 通过反射TypeOf/ValueOf/Field/NumMethod等方法获取接口对象的字段,类型和方法等信息

package main

import (
	"fmt"
	"reflect"
)

type TUser struct {
	Id   int
	Name string
	Age  int
}

func (u TUser) Hello() {
	fmt.Println("Hello,Mike")
}

func StructReflectInit() {
	u := TUser{1, "Mike", 11}
	Info(u)
}

func Info(o interface{}) {
	t := reflect.TypeOf(o)         //获取接口的类型
	fmt.Println("Type:", t.Name()) //t.Name() 获取接口的名称

	if t.Kind() != reflect.Struct {
		fmt.Println("err: type invalid")
		return
	}

	v := reflect.ValueOf(o) // 获取接口的值类型
	fmt.Println("Fields: ")

	for i := 0; i < t.NumField(); i++ { // NumField取出这个接口所有的字段数量
		f := t.Field(i)                                   // 获得结构体的第i个字段
		val := v.Field(i).Interface()                     // 获得接口子弹的值
		fmt.Printf("%6s: %v = %v\n", f.Name, f.Type, val) //第i个字段第名称,类型,值
	}

	for i := 0; i < t.NumMethod(); i++ { // NumMethod 获得这个接口所有的方法数量
		m := t.Method(i)
		fmt.Printf("%6s: %v\n", m.Name, m.Type) // 输出方法名称和方法类型
	}

}

func main() {
	StructReflectInit()
}
