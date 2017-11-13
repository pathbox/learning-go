package main

import(
	"fmt"

    "./protobuf"
    "github.com/golang/protobuf/proto"
)

func main(){
	data := protobuf.User{
		Id: 12345,
		Username: "John Harry",
		Password: "password",
	}

	// 将资料编码成 Protocol Buffer 格式 传入的是 Pointer
	dataBuffer, _ := proto.Marshal(&data)

	// 将已经编码的资料解码成 protobuf.User 格式
	var user protobuf.User
	proto.Unmarshal(dataBuffer, &user)

	fmt.Println(user.Id, " ", user.Username, " ", user.Password)
}

/*

过程梳理

1. 准备user.proto格式文件
2. protoc --go_out=. user.proto 会得到user.pb.go 文件，将其放到protobuf文件夹中 在user.proto中定义了package为 protobuf，这个可以自行定义
这样，生成的user.pb.go 为package protobuf
3. 使用。在main.go 导入protobuf package。
导入 protoc --go_out=. user.proto

protobuf.User{}来构建 struct数据

序列化
proto.Marshal(&data) 序列化struct，之后将dataBuffer 进行网络传递

反序列化
var user protobuf.User 定义userstruct
proto.Unmarshal(dataBuffer, &user) 反序列化数据到user struct

*/
