package main

import (
    pb "../../proto/hello" // 引入编译生成的包

		"golang.org/x/net/context"  // 不能使用  context包...
    "google.golang.org/grpc"
		"google.golang.org/grpc/grpclog"
		"fmt"
)
const (
    // Address gRPC服务地址
    Address = "127.0.0.1:50052"
)

func main() {
	// 连接
	conn, err := grpc.Dial(Address, grpc.WithInsecure())

	if err != nil {
		grpclog.Fatalln(err)
	}

	defer conn.Close()

	// 初始化客户端, conn 作为参数
	client := pb.NewHelloClient(conn)

	// 调用方法
	reqBody := new(pb.HelloRequest)
	reqBody.Name = "gRPC"
	r, err := client.SayHello(context.Background(), reqBody) // 远程调用方法
	if err != nil {
    grpclog.Fatalln(err)
  }

  fmt.Println("Message: ", r.Message)
}