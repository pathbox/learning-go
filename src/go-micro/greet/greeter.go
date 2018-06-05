package main

import (
	"context"
	"log"

	"github.com/micro/go-micro"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	rsp.Greeting = "Hello " + req.Name
	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("greeter"),
		micro.Version("lastest"),
	)

	service.Init()
	proto.RegisterGreeterHandler(service.Server(), new(Greeter)) // 注册handler
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
