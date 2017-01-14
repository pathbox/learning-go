package main

import (
	"log"
	"net/rpc"

	"./rpcexample"
)

func main() {
	client, err := rpc.DialHTTP("tcp", ":1234")
	if err != nil {
		log.Fatalf("Error in dialing. %s", err)
	}

	args := &rpcexample.Args{
		A: 2,
		B: 3,
	}

	var result rpcexample.Result

	err = client.Call("Arith.Multiply", args, &result)
	if err != nil {
		log.Fatalln("error in Arith", err)
	}
	log.Printf("%d*%d=%d", args.A, args.B, result)
}
