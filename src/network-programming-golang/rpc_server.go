package main

import (
	"fmt"

	"net/rpc"

	"errors"

	"net/http"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {

	*reply = args.A * args.B

	return nil

}

func (t *Arith) Divide(args *Args, quo *Quotient) error {

	if args.B == 0 {

		return errors.New("divide by zero")

	}

	quo.Quo = args.A / args.B

	quo.Rem = args.A % args.B

	return nil

}

func main() {

	arith := new(Arith)

	rpc.Register(arith)

	rpc.HandleHTTP()

	err := http.ListenAndServe(":12345", nil)

	if err != nil {

		fmt.Println(err.Error())

	}

}
