package main

import (
	"encoding/json"
	"fmt"
)

// Decode
// 将JSON数据解码

// func Unmarshal(data []byte, v interface{}) error

type Animal struct {
	Name  string
	Order string
}

func main() {
	var jsonBlob = []byte(`[
    {"Name": "Platypus", "Order": "Monotremata"},
    {"Name": "Quoll",    "Order": "Dasyuromorphia"}
]`)
	var animals []Animal
	err := json.Unmarshal(jsonBlob, &animals)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", animals)
}
