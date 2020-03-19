package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Policy struct {
	Statement []Statement `json:"Statement"`
}

type Statement struct {
	Effect    string                 `json:"Effect"`
	Action    interface{}            `json:"Action"`
	Resource  interface{}            `json:"Resource"`
	Condition map[string]interface{} `json:"Condition,omitempty"`
}

func main() {
	// // 声明一个空结构体
	// type cat struct {
	// }
	// // 获取结构体实例的反射类型对象
	// typeOfCat := reflect.TypeOf(cat{})
	// // 显示反射类型对象的名称和种类
	// fmt.Println(typeOfCat.Name(), typeOfCat.Kind())
	// // 获取Zero常量的反射类型对象
	// typeOfA := reflect.TypeOf(Zero)
	// // 显示反射类型对象的名称和种类
	// fmt.Println(typeOfA.Name(), typeOfA.Kind())
	jsonByte := []byte(`{
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                      "oss:ListBuckets",
                      "oss:GetBucketStat",
                      "oss:GetBucketInfo",
                      "oss:GetBucketTagging",
                      "oss:GetBucketAcl"
                      ],
            "Resource": [
                "acs:oss:*:*:*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "oss:GetObject",
                "oss:GetObjectAcl"
            ],
            "Resource": [
                "acs:oss:*:*:myphotos/hangzhou/2015/*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": "oss:ListObjects",
            "Resource": [
                "acs:oss:*:*:myphotos"
            ],
            "Condition": {
                "StringLike": {
                    "oss:Delimiter": "/",
                    "oss:Prefix": [
                        "",
                        "hangzhou/",
                        "hangzhou/2015/*"
                    ]
                }
            }
        }
    ]
}`)

	py := &Policy{}
	err := json.Unmarshal(jsonByte, &py)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("Statement: %+v\n", py)

	for _, st := range py.Statement {
		action := st.Action
		typeOfA := reflect.TypeOf(action)
		valueOfA := reflect.ValueOf(action)
		// 显示反射类型对象的名称和种类
		// fmt.Println("===", typeOfA.Name(), typeOfA.Kind(), valueOfA)
		var r string
		switch typeOfA.Kind() {
		case reflect.Slice:
			// te := typeOfA.Elem()
			for i := 0; i < valueOfA.Len(); i++ {
				val := valueOfA.Index(i)
				fmt.Println("tttt:", val.Kind()) // 是interface 类型
				vala := val.Interface().(string)
				fmt.Println("vvvvvv:", vala)
				r = fmt.Sprintf("%s", val)
				fmt.Println("rrrrrr:", r)
			}

		case reflect.String:
			fmt.Printf("Action: %s\n", valueOfA.String())
		}
	}
}
