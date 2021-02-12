package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/selector"
	"github.com/asim/go-micro/v3/cmd"
	"github.com/asim/go-micro/v3/metadata"
	"github.com/asim/go-micro/v3/registry"

	example "github.com/asim/go-micro/examples/v3/server/proto/example"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type dcWrapper struct {
	client.Client
}

// 包裹了client.Client，然后实现Call方法
func (dc *dcWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	md,_ := metadata.FromContext(ctx)

	filter := func(services []*registry.Service) []*registry.Service {
		for _, service := range services {
			var nodes []*registry.Node
			for _, node := range service.Nodes {
				if node.Metadata["datacenter"] == md["datacenter"] {
					nodes = append(nodes, node)
				}
			}
				service.Nodes = nodes
		}
		return services
	}
// 定义callOptions，传入func filter，满足node.Metadata["datacenter"] == md["datacenter"] 请求则到相应的node中
	callOptions := append(opts, client.WithSelectOption(
		selector.WithFilter(filter),
	))

	fmt.Printf("[DC Wrapper] filtering for datacenter %s\n", md["datacenter"])
	return dc.Client.Call(ctx, req, rsp, callOptions...)
}

func NewDCWrapper(c client.Client) client.Client {
	return &dcWrapper{c}
}

func call(i int) {
	req := client.NewRequest("go.micro.srv.example", "Example.Call",&example.Request{
		Name: "John",
	})


	ctx := metadata.NewContext(context.Background(), map[string]string{
		"datacenter": "local",
	})

	if err := client.Call(ctx, req, rsp); err != nil {
		fmt.Println("call err: ", err, rsp)
		return
	}
	fmt.Println("Call:", i, "rsp:", rsp.Msg)
}

func main() {
	cmd.Init()

	client.DefaultClient = client.NewClient(
		client.Wrap(NewDCWrapper)
	)

	fmt.Println("\n--- Call example ---")
	for i := 0; i < 10; i++ {
		call(i)
	}
}