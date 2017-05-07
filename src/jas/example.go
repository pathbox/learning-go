// A simple and powerful REST API framework for Go

package main

import (
	"github.com/coocood/jas"
	"net/http"
)

type Hello struct{}

func (*Hello) Get(ctx *jas.Context) {
	ctx.Data = "hello world"
	// response: `{"data": "hello world", "error": null}`
}

func main() {
	router := jas.NewRouter(new(Hello))
	router.BasePath = "/v1/"
	fmt.Println(router.HandledPaths(true))
	//output: `GET /v1/hello`
	http.Handle(router.BasePath, router)
	http.ListenAndServe(":8080", nil)

}
