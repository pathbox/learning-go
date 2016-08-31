package main


import (
	"fmt"
	"reflect"
)

// 计算器程序 函数按名称调用

/*
最主要的是根据函数所在的结构体对象根据函数名称获取函数值，
即通过reflect.TypeOf(calc).MethodByName(funcName)
来获取函数的值，然后再使用[]reflect.Value来组装函数调用所需要的参数，
其中最重要的是params[0]必须是函数所在结构体对象的值，其他的则为对应函数的参数了。
最后使用fn.Func.Call方法调用函数即可。
*/

var opFuncs = map[string]string{
	"+": "Add",
	"-": "Minus",
	"*": "Multi",
	"/": "Divide",
}

type Calculator struct {
}

func (this *Calculator) Add(a, b int) int {
  return a + b
}

func (this *Calculator) Minus(a, b int) int {
  return a - b
}

func (this *Calculator) Multi(a, b int) int {
  return a * b
}

func (this *Calculator) Divide(a, b int) int {
  return a / b
}

func main() {
  var a int = 10
  var b int = 20

  calc := &Calculator{}
  for funcOp, funcName := range opFuncs {
    fn, _ := reflect.TypeOf(calc).MethodByName(funcName)

    params := make([]reflect.Value, 3)
    params[0] = reflect.ValueOf(calc)
		params[1] = reflect.ValueOf(a)
		params[2] = reflect.ValueOf(b)
    v := fn.Func.Call(params)

    fmt.Println(fmt.Sprintf("%d %s %d = %d", a, funcOp, b, v[0].Int()))
  }
}
