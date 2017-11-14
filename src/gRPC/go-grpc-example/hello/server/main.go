// protoc -I . --go_out=plugins=grpc:. ./hello.proto

package main

import (
    "net"

    pb "../../proto/hello" // 引入编译生成的包

		"golang.org/x/net/context"  // 不能使用  context包...
    "google.golang.org/grpc"
    "google.golang.org/grpc/grpclog"
)

const (
	// Address gRPC 服务地址
	Address = "127.0.0.1:50052"
)

// 定义helloService并实现约定的接口
type helloService struct{}

func (h helloService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	resp := new(pb.HelloReply)
	resp.Message = "Hello " + in.Name + "."
	return resp, nil
}

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("falied to listen: %v", err)
	}
	// 实例化 grpc Server
	s := grpc.NewServer()

	// 注册HelloService
	pb.RegisterHelloServer(s, &helloService{})

	grpclog.Println("Listen on " + Address)

	s.Serve(listen)

}