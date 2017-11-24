package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"../protobuf"
)

// server 结构体会作为 Calculator的gRPC服务器
type server struct{}

func (s *server) Plus(ctx context.Context, in *protobuf.CalcRequest) (*protobuf.CalcReply, error) {
	result := in.NumberA + in.NumberB

	// 包装成Protobuf 结构体并返回
	return &protobuf.CalcReply{Result: result}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}

	// 建立新 gRPC服务器并注册Calculator服务
	s := grpc.NewServer()                           // 新建一个grpc服务
	protobuf.RegisterCalculatorServer(s, &server{}) // 将server 和 s 绑定，server的plus方法会在gRPC中调用

	// 在gRPC服务器上注册反射服务
	reflection.Register(s)

	// 开始在指定的端口中服务
	if err := s.Serve(listen); err != nil {
		log.Fatalln(err)
	}
}
