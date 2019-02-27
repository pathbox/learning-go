package main

import (
	"io"
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

type bidirectionalStreamServer struct{}

// bidirectionalStreamServer 实现了 user.pb.go 中的 UserServiceServer 接口
func (s *bidirectionalStreamServer) GetUserInfo(stream pb.UserService_GetUserInfoServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		u := users[req.ID]
		err = stream.Send(&u)
		if err != nil {
			return err
		}
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
	pb.RegisterUserServiceServer(grpcServer, &bidirectionalStreamServer{})

	// 启动 gRPC 服务器
	// 阻塞等待客户端的调用
	grpcServer.Serve(listener)
}
