package main

import (
	"context"
	"log"
	"net"

	pb "../proto"
	"google.golang.org/grpc"
)

// 模拟的数据库查询结果
var users = map[int32]pb.UserResponse{
	1: {Name: "Dennis MacAlistair Ritchie", Age: 70},
	2: {Name: "Ken Thompson", Age: 75},
	3: {Name: "Rob Pike", Age: 62},
}

type simpleServer struct{}

// simpleServer 实现了 user.pb.go 中的 UserServiceServer 接口
func (s *simpleServer) GetUserInfo(ctx context.Context, req *pb.UserRequest) (resp *pb.UserResponse, err error) {
	if user, ok := users[req.ID]; ok {
		resp = &user
	}
	log.Printf("[RECEVIED REQUEST]: %v\n", req)
	return
}

func main() {
	// 1.指定微服务的服务端监听地址
	addr := "0.0.0.0:9999"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	} else {
		log.Println("server listen: ", addr)
	}

	// 2.创建 gRPC 服务器实例
	grpcServer := grpc.NewServer()

	// 3.向gRPC服务注册服务
	pb.RegisterUserServiceServer(grpcServer, &simpleServer{})

	// 4. 启动 gRPC 服务器
	// 阻塞等待客户端的调用
	grpcServer.Serve(listener)
}
