package main

import (
	"fmt"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

// fasthttprouter.RequestCtx.UserValue() 可以获得路由匹配得到的参数，如规则 /hello/:name 中的 :name
func httpHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hello, %s!\n", ctx.UserValue("name"))
}

func main() {
	//　使用　fasthttprouter 创建路由
	router := fasthttprouter.New()
	router.GET("/hello/:name", httpHandler)
	log.Println(fasthttp.ListenAndServe(":9090", router.Handler))
}
