package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Id      int
	Name    string
	Address Address
}

type Address struct {
	Add string
	Res int
}

func main() {
	u := User{Id: 1001, Name: "aaa", Address: Address{Add: "ccccccccc", Res: 12}}
	t := reflect.TypeOf(u)
	v := reflect.ValueOf(u)

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() { // 是否为可导出字段
			// 判断是否是嵌套结构
			if v.Field(i).Type().Kind() == reflect.Struct {
				fmt.Println("----------------")
				structField := v.Field(i).Type()
				for j := 0; j < structField.NumField(); j++ {
					fmt.Printf("%s %s = %v -tag:%s \n",
						structField.Field(j).Name,
						structField.Field(j).Type,
						v.Field(i).Field(j).Interface(),
						structField.Field(j).Tag)
				}
				continue
			}
			fmt.Printf("%s %s = %v -tag:%s \n",
				t.Field(i).Name,
				t.Field(i).Type,
				v.Field(i).Interface(),
				t.Field(i).Tag)
		}
	}
}
