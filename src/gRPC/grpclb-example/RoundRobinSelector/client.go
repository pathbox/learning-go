package main

import (
	"log"

	etcd "github.com/coreos/etcd/client"
	grpclb "github.com/liyue201/grpc-lb"
	"github.com/liyue201/grpc-lb/examples/proto"
	registry "github.com/liyue201/grpc-lb/registry/etcd"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	etcdConfg := etcd.Config{
		Endpoints: []string{"http://127.0.0.1:4001"},
	}
	r := registry.NewResolver("/grpc-lb", "test", etcdConfg)
	b := grpclb.NewBalancer(r, grpclb.NewRoundRobinSelector())
	c, err := grpc.Dial("", grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		log.Printf("grpc dial: %s", err)
		return
	}
	defer c.Close()
	client := proto.NewTestClient(c)
	resp, err := client.Say(context.Background(), &proto.SayReq{Content: "round robin"})
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf(resp.Content)

}
