package main

import (
	"net/rpc"

	"fmt"

	"log"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

func main() {

	serverAddress := "127.0.0.1:12345"

	client, err := rpc.DialHTTP("tcp", serverAddress)

	if err != nil {

		log.Fatal("dialing:", err)

	}

	// Synchronous call

	args := Args{17, 8}

	var reply int

	err = client.Call("Arith.Multiply", args, &reply)

	if err != nil {

		log.Fatal("arith error:", err)

	}

	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

	var quot Quotient

	err = client.Call("Arith.Divide", args, &quot)

	if err != nil {

		log.Fatal("arith error:", err)

	}

	fmt.Printf("Arith: %d/%d=%d remainder %d\n", args.A, args.B, quot.Quo, quot.Rem)

}
