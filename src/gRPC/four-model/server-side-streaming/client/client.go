package main

import (
	"context"
	"io"
	"log"

	pb "../proto"

	"google.golang.org/grpc"
)

func main() {
	// gRPC 服务器的地址
	addr := "0.0.0.0:9999"

	// 不使用认证建立连接
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect server error: %v", err)
	}
	defer conn.Close()

	// 创建 gRPC 客户端实例
	grpcClient := pb.NewUserServiceClient(conn)

	// 调用服务端的函数
	req := pb.UserRequest{ID: 1}
	stream, err := grpcClient.GetUserInfo(context.Background(), &req)
	if err != nil {
		log.Fatalf("recevie resp error: %v", err)
	}

	// 接收流数据
	for {
		resp, err := stream.Recv()
		if err == io.EOF { // 服务端数据发送完毕
			break
		}
		if err != nil {
			log.Fatalf("receive error: %v", err)
		}
		log.Printf("[RECEIVED RESPONSE]: %v\n", resp) // 输出响应
	}
}
