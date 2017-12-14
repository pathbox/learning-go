package main

import (
	"encoding/json"
	"fmt"
)

type Resp struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

func main() {
	data := "[{\"name\":\"John\", \"id\": 1}]"

	var resp []Resp
	json.Unmarshal([]byte(data), &resp)

	fmt.Printf("%+v", resp)

}

// 序列化 数组字符串为数组，数组中的元素为 json对象。struct 中 field 首字母要大写
