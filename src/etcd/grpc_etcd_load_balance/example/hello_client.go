package main

import (
	"flag"
	"fmt"
	"time"

	"golang.org/x/net/context"

	balance "../balance"
	pb "./pb"
	"google.golang.org/grpc"
)

var (
	serv = flag.String("service", "hello_service", "service name")
	reg  = flag.String("reg", "http://127.0.0.1:2379", "register etcd address")
)

func main() {
	flag.Parse()
	fmt.Println("serv", *serv)
	resolver := balance.NewResolver(*serv) // 生成一个resolver
	b := grpc.RoundRobin(resolver)         // 使用轮询负载

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	conn, err := grpc.DialContext(ctx, *reg, grpc.WithInsecure(), grpc.WithBalancer(b)) // 使用上面的配置,Dial一个conn
	if err != nil {
		panic(err)
	}

	fmt.Println("conn...")

	client := pb.NewGreeterClient(conn) // 使用该conn,生成client
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "world"})
	if err == nil {
		fmt.Printf("Reply is %s\n", resp.Message)
	} else {
		fmt.Println(err)
	}

}
