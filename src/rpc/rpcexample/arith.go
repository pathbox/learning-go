package rpcexample

import (
	"log"
)

//Holds arguments to be passed to service Arith in RPC call
type Args struct {
	A, B int
}

//Representss service Arith with method Multiply
type Arith int

//Result of RPC call is of this type
type Result int

//This procedure is invoked by rpc and calls rpcexample.Multiply which stores product of args.A and args.B in result pointer
func (t *Arith) Multiply(args Args, result *Result) error {
	return Multiply(args, result)
}

//stores product of args.A and args.B in result pointer
func Multiply(args Args, result *Result) error {
	log.Printf("Multiplying %d with %d\n", args.A, args.B)
	*result = Result(args.A * args.B)
	return nil
}
