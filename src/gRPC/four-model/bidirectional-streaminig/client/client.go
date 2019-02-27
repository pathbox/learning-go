package main

import (
	"context"
	"log"
	"time"

	pb "../proto"
	"google.golang.org/grpc"
)

func main() {
	// gRPC 服务器的地址
	addr := "0.0.0.0:2333"

	// 不使用认证建立连接
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect server error: %v", err)
	}
	defer conn.Close()

	// 创建 gRPC 客户端实例
	grpcClient := pb.NewUserServiceClient(conn)
	stream, err := grpcClient.GetUserInfo(context.Background())
	if err != nil {
		log.Fatalf("receive stream error: %v", err)
	}

	// 向服务端发送数据流，并处理响应流
	var i int32
	for i = 1; i < 4; i++ {
		stream.Send(&pb.UserRequest{ID: i})
		time.Sleep(1 * time.Second)
		resp, err := stream.Recv()
		if err != nil {
			log.Fatalf("resp error: %v", err)
		}
		log.Printf("[RECEIVED RESPONSE]: %v\n", resp) // 输出响应
	}
}
