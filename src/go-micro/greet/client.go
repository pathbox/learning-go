package main

import (
	"context"
	"fmt"
)

func main() {
	// create the greeter client using the service name and client
	greeter := proto.NewGreeterClient("greeter", service.Client())
	// request the Hello method on the Greeter handler
	rsp, err := greeter.Hello(context.TODO(), &proto.HelloRequest{
		Name: "John",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp.Greeter)
}
