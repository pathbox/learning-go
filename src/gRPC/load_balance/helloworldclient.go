package main

import (
    "flag"
    "fmt"
    "time"

    grpclb "com.midea/jr/grpclb/naming/etcd/v3"
    "com.midea/jr/grpclb/example/pb"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
    "strconv"
)

var (
    serv = flag.String("service", "hello_service", "service name")
    reg = flag.String("reg", "http://127.0.0.1:2379", "register etcd address")
)

func main() {
    flag.Parse()
    r := grpclb.NewResolver(*serv)
    b := grpc.RoundRobin(r)  // 在client端决定使用哪种轮询算法

    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    conn, err := grpc.DialContext(ctx, *reg, grpc.WithInsecure(), grpc.WithBalancer(b))
    if err != nil {
        panic(err)
    }

    ticker := time.NewTicker(1 * time.Second)
    for t := range ticker.C {
        client := pb.NewGreeterClient(conn)
        resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "world " + strconv.Itoa(t.Second())})
        if err == nil {
            fmt.Printf("%v: Reply is %s\n", t, resp.Message)
        }
    }
}