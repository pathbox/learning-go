package json

import (
	"net/http"

	"github.com/haisum/rpcexample"
)

//Represents service Arith on JSON-RPC
type Arith int

//Invoked by JSON-RPC client and calls rpcexample.Multiply which stores product of args.A and args.B in result
func (t *Arith) Multiply(r *http.Request, args *rpcexample.Args, result *rpcexample.Result) error {
	return rpcexample.Multiply(*args, result)
}
