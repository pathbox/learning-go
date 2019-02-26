package main

import (
	"context"
	"log"

	pb "../proto"

	"google.golang.org/grpc"
)

func main() {
	addr := "0.0.0.0:9999"

	// 1. Dial grpc conn
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect server error: %v", err)
	}
	defer conn.Close()

	// 2. 创建 gRPC 客户端实例
	grpcClient := pb.NewUserServiceClient(conn)

	// 3. 准备请求结构

	req := pb.UserRequest{ID: 2}
	// 4. client 调用服务端实现的接口方法
	resp, err := grpcClient.GetUserInfo(context.Background(), &req)
	if err != nil {
		log.Fatalf("recevie resp error: %v", err)
	}

	log.Printf("[RECEIVED RESPONSE]: %v\n", resp)
}
