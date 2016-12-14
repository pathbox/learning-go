package main

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Name  string
	Body  string
	Time  int64
	inner string
}

func main() {
	var m = Message{
		Name:  "Alice",
		Body:  "Hello",
		Time:  1294706395881547000,
		inner: "ok",
	}
	b := []byte(`{"Name":"Bob","Food":"Pickle", "inner":"changed"}`)

	err := json.Unmarshal(b, &m)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	fmt.Printf("%v", m)
}

// 结构体必须是大写字母开头的成员才会被JSON处理到，小写字母开头的成员不会有影响。

// Mashal时，结构体的成员变量名将会直接作为JSON Object的key打包成JSON；Unmashal时，会自动匹配对应的变量名进行赋值，大小写不敏感。
// Unmarshal时，如果JSON中有多余的字段，会被直接抛弃掉；如果JSON缺少某个字段，则直接忽略不对结构体中变量赋值，不会报错

// 更灵活地使用JSON

// 使用json.RawMessage
// json.RawMessage其实就是[]byte类型的重定义。可以进行强制类型转换。

// 现在有这么一种场景，结构体中的其中一个字段的格式是未知的：

// [plain] view plain copy 在CODE上查看代码片派生到我的代码片
// type Command struct {
//     ID   int
//     Cmd  string
//     Args *json.RawMessage
// }

// 使用interface{}
// interface{}类型在Unmarshal时，会自动将JSON转换为对应的数据类型：

// JSON的boolean 转换为bool
// JSON的数值 转换为float64
// JSON的字符串 转换为string
// JSON的Array 转换为[]interface{}
// JSON的Object 转换为map[string]interface{}
// JSON的null 转换为nil

// 需要注意的有两个。一个是所有的JSON数值自动转换为float64类型，使用时需要再手动转换为需要的int，int64等类型。第二个是JSON的object自动转换为map[string]interface{}类型，访问时直接用JSON Object的字段名作为key进行访问。再不知道JSON数据的格式时，可以使用interface{}。

// 自定义类型
// 如果希望自己定义对象的打包解包方式，可以实现以下的接口：

// [plain] view plain copy 在CODE上查看代码片派生到我的代码片
// type Marshaler interface {
//     MarshalJSON() ([]byte, error)
// }
// type Unmarshaler interface {
//     UnmarshalJSON([]byte) error
// }
// 实现该接口的对象需要将自己的数据打包和解包。如果实现了该接口，json在打包解包时则会调用自定义的方法，不再对该对象进行其他处理。
