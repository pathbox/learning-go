package main

import (
	"context"
	pb "grpc/client-side-streaming/proto"
	"log"

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

	// 向服务端发送流数据
	stream, err := grpcClient.GetUserInfo(context.Background())

	var i int32
	// 模拟的数据库中有 3 条记录，ID 分别为 1 2 3
	for i = 1; i < 4; i++ {
		err := stream.Send(&pb.UserRequest{ID: i})
		if err != nil {
			log.Fatalf("send error: %v", err)
		}
	}

	// 接收服务端的响应
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("recevie resp error: %v", err)
	}

	log.Printf("[RECEIVED RESPONSE]: %v\n", resp) // 输出响应
}
