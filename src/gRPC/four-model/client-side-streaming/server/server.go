package main

import (
	pb "grpc/client-side-streaming/proto"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

var users = map[int32]pb.UserResponse{
	1: {Name: "Dennis MacAlistair Ritchie", Age: 70},
	2: {Name: "Ken Thompson", Age: 75},
	3: {Name: "Rob Pike", Age: 62},
}

type clientSideStreamServer struct{}

func (s *clientSideStreamServer) GetUserInfo(stream pb.UserService_GetUserInfoServer) error {
	var lastID int32
	for {
		req, err := stream.Recv()
		// 客户端数据流发送完毕
		if err == io.EOF {
			// 返回最后一个 ID 的用户信息
			if u, ok := users[lastID]; ok {
				stream.SendAndClose(&u)
				return nil
			}
		}
		lastID = req.ID
		log.Printf("[RECEVIED REQUEST]: %v\n", req)
	}
	return nil
}

func main() {
	// 指定微服务的服务端监听地址
	addr := "0.0.0.0:2333"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	} else {
		log.Println("server listen: ", addr)
	}

	// 创建 gRPC 服务器实例
	grpcServer := grpc.NewServer()

	// 向 gRPC 服务器注册服务
	pb.RegisterUserServiceServer(grpcServer, &clientSideStreamServer{})

	// 启动 gRPC 服务器
	// 阻塞等待客户端的调用
	grpcServer.Serve(listener)
}
