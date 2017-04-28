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
	ctx.SetContentType("text/html") // 如果不被添加这个，　都当纯文本返回
	fmt.Fprintf(ctx, "Method: %s <br/>", ctx.Method())
	fmt.Fprintf(ctx, "URI: %s <br/>", ctx.URI())
	fmt.Fprintf(ctx, "RemoteAddr: %s <br/>", ctx.RemoteAddr())
	fmt.Fprintf(ctx, "UserAgent: %s <br/>", ctx.UserAgent())
	fmt.Fprintf(ctx, "Header.Accept: %s <br/>", ctx.Request.Header.Peek("Accept"))
	fmt.Fprintf(ctx, "IP:%s <br/>", ctx.RemoteIP())
	fmt.Fprintf(ctx, "Host:%s <br/>", ctx.Host())
	fmt.Fprintf(ctx, "ConnectTime:%s <br/>", ctx.ConnTime()) // 连接收到处理的时间
	fmt.Fprintf(ctx, "IsGET:%v <br/>", ctx.IsGet())          // 类似有 IsPOST, IsPUT 等

}

func main() {
	//　使用　fasthttprouter 创建路由
	router := fasthttprouter.New()
	router.GET("/hello/:name", httpHandler)
	log.Println(fasthttp.ListenAndServe(":9090", router.Handler))
}
