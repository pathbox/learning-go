package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type Person struct {
	Name   string
	Age    int
	Avatar struct {
		Url    string
		Height int
		Width  int
	}
}

func main() {
	fmt.Println(MyStruct(reflect.TypeOf(Person{})))
}

func MyStruct(t reflect.Type) reflect.Type {
	if t.Kind() != reflect.Struct {
		panic("invalid type")
	}

	fs := make([]reflect.StructField, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fs[i] = f
		if f.Tag.Get("json") == "" {
			fs[i].Tag = reflect.StructTag(`json:"` + strings.ToLower(f.Name) + `"`)
		}
		var fType reflect.Type
		switch f.Type.Kind() {
		case reflect.Struct:
			fType = MyStruct(f.Type)
		case reflect.Slice:
			if f.Type.Elem().Kind() == reflect.Struct {
				fType = reflect.SliceOf(MyStruct(f.Type.Elem()))
			} else if f.Type.Elem().Kind() == reflect.Slice {
				panic("multi-d slice not supported") //多维数组暂不考虑
			}
		default:
			fType = f.Type
		}
		fs[i].Type = fType
	}
	fmt.Printf("%+v\n", fs)
	fmt.Println()
	return reflect.StructOf(fs)
}

func MyMarshal(obj interface{}) (b []byte, e error) {
	b, e = json.Marshal(obj)
	if e != nil {
		return
	}
	var m map[string]interface{}
	e = json.Unmarshal(b, &m)
	if e != nil {
		return
	}
	// 序列化反序列化之后再传入m
	HandleMapStyle(m)
	return json.Marshal(m)
}

func HandleMapStyle(m map[string]interface{}) {
	for key, value := range m {
		switch v := value.(type) {
		case []interface{}:
			for _, i := range v {
				if elem, ok := i.(map[string]interface{}); ok {
					HandleMapStyle(elem)
				}
			}
		case map[string]interface{}:
			HandleMapStyle(v)
		}
		delete(m, key) // 如果是 []interface{} 或 map 就畸形递归处理，如果不是，则直接进行下面
		m[strings.ToLower(key)] = value
	}
}
