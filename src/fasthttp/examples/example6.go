package main

import (
	"fmt"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func httpHandler(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("hello fasthttp")
	// 因为实现不同，fasthttp 的返回内容不是即刻返回的
	// 不同于标准库，添加返回内容后设置状态码，也是有效的
	ctx.SetStatusCode(404)

	// 返回的内容也是可以获取的，不需要标准库的用法，需要自己扩展 http.ResponseWriter
	fmt.Println(string(ctx.Response.Body()))
}

func main() {
	//　使用　fasthttprouter 创建路由
	router := fasthttprouter.New()
	router.GET("/", httpHandler)
	log.Println(fasthttp.ListenAndServe(":9090", router.Handler))
}
