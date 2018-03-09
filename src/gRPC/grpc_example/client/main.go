package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"./../protobuf"
)

func main() {
	// 连接到远端gRPC服务器
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()

	// 建立新的Calculator客户端 所以等一下就能够使用Calculator的所有方法
	c := protobuf.NewCalculatorClient(conn)

	// 传送新请求到远端gRPC服务器Calculator中，并呼叫Plus函数
	r, err := c.Plus(context.Background(), &protobuf.CalcRequest{NumberA: 32, NumberB: 88})
	if err != nil {
		log.Panicln(err)
	}

	log.Println("得到的结果是： ", r.Result)
}

// 你可以直接在本地程式直接呼叫遠端程式，用法十分貼切。因為我們有 Protobuf 協助我們在客戶端和伺服端都共享同一份協定、結構
