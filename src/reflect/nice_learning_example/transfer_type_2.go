// 未知原有类型【遍历探测其Filed】
// 很多情况下，我们可能并不知道其具体类型，那么这个时候，该如何做呢？需要我们进行遍历探测其Filed来得知

package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Id   int
	Name string
	Age  int
}

func (u User) ReflectCallFunc() {
	fmt.Println("Cool ReflectCallFunc")
}

func main() {
	user := User{1, "Cool Guy", 27}
	DoFiledAndMethod(user)
}

// 通过接口来获取任意参数，然后一一揭晓
func DoFiledAndMethod(input interface{}) {
	getType := reflect.TypeOf(input)
	fmt.Println("get Type is :", getType.Name())

	getValue := reflect.ValueOf(input)
	fmt.Println("get all Fields is:", getValue)

	// 获取方法字段
	// 1. 先获取interface的reflect.Type，然后通过NumField进行遍历
	// 2. 再通过reflect.Type的Field获取其Field
	// 3. 最后通过Field的Interface()得到对应的value
	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		// field := getType.FieldByIndex([]int{i})
		value := getValue.Field(i).Interface() // 不再继续强制类型转换
		fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
	}

	// 获取方法
	// 先获取interface的reflect.Type, 然后通过.NumMethod进行遍历
	for i := 0; i < getType.NumMethod(); i++ {
		m := getType.Method(i)
		fmt.Printf("%s: %v\n", m.Name, m.Type)
	}

}

/*
通过运行结果可以得知获取未知类型的interface的具体变量及其类型的步骤为：

先获取interface的reflect.Type，然后通过NumField进行遍历
再通过reflect.Type的Field获取其Field
最后通过Field的Interface()得到对应的value
通过运行结果可以得知获取未知类型的interface的所属方法（函数）的步骤为：

先获取interface的reflect.Type，然后通过NumMethod进行遍历
再分别通过reflect.Type的Method获取对应的真实的方法（函数）
最后对结果取其Name和Type得知具体的方法名
也就是说反射可以将“反射类型对象”再重新转换为“接口类型变量”
struct 或者 struct 的嵌套都是一样的判断处理方式
*/
