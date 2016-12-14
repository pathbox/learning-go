package main

import (
	"encoding/json"
	"os"
)

type ColorGroup struct {
	ID     int
	Name   string
	Colors []string
}

func main() {
	group := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}

	json, err := json.Marshal(group)
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(json)
}

// Encode
// 将一个对象编码成JSON数据，接受一个interface{}对象，返回[]byte和error：
// func Marshal(v interface{}) ([]byte, error)
// Marshal函数将会递归遍历整个对象，依次按成员类型对这个对象进行编码，类型转换规则如下：

// bool类型 转换为JSON的Boolean
// 整数，浮点数等数值类型 转换为JSON的Number
// string 转换为JSON的字符串(带""引号)
// struct 转换为JSON的Object，再根据各个成员的类型递归打包
// 数组或切片 转换为JSON的Array
// []byte 会先进行base64编码然后转换为JSON字符串
// map 转换为JSON的Object，key必须是string
// interface{} 按照内部的实际类型进行转换
// nil 转为JSON的null
// channel,func等类型 会返回UnsupportedTypeError
