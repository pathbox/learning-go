// 通过反射“动态”调用方法

package main

import (
	"fmt"
	"reflect"
)

type FUser struct {
	Id   int
	Name string
	Age  int
}

func (u FUser) Hello(m FUser) (int, string) {
	fmt.Println("Hello", m.Name, ", I am ", u.Name)
	return m.Age + u.Age, u.Name
}

func FuncReflectInit() {
	u := FUser{1, "Mike", 25}
	GetInfo(u)
}

func GetInfo(u interface{}) {
	m := FUser{2, "Json", 29}

	v := reflect.ValueOf(u)
	if v.Kind() != reflect.Struct {
		fmt.Println("type invalid")
		return
	}

	mv := v.MethodByName("Hello") // 获取对应的方法
	if !mv.IsValid() {            //判断方法是否存在
		fmt.Println("Func Hello not exist")
		return
	}

	args := []reflect.Value{reflect.ValueOf(m)} //初始化传入等参数，传入等类型只能是[]reflect.Value类型
	res := mv.Call(args)
	fmt.Println(res[0], res[1])
}

func main() {
	FuncReflectInit()
}

// 通过MethodByName先获取对象的Hello方法,然后准备要传入的参数,这里传入的参数必须是[]refelct.Value类型,传入的参数值必须强制转换为反射值类型refelct.Value。
// 最后通过调用Call方法就可以实现通过反射”动态”调用对象的方法
