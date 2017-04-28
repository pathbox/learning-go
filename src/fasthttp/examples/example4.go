package main

import (
	"fmt"
	"log"

	"bytes"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

// fasthttprouter.RequestCtx.UserValue() 可以获得路由匹配得到的参数，如规则 /hello/:name 中的 :name
func httpHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hello, %s!\n", ctx.UserValue("name"))
	ctx.SetContentType("text/html") // 如果不被添加这个，　都当纯文本返回

	// GET ?abc=abc&abc=123
	getValues := ctx.QueryArgs()
	fmt.Fprintf(ctx, "GET abc=%s <br/>", getValues.Peek("abc")) // Peek 只获取第一个值
	fmt.Fprintf(ctx, "GET abc=%s <br/>",
		bytes.Join(getValues.PeekMulti("abc"), []byte(","))) // PeekMulti 获取所有值

	// POST xyz = xyz & xyz = 123
	postValues := ctx.PostArgs()
	fmt.Fprintf(ctx, "POST xyz=%s <br/>",
		postValues.Peek("xyz"))
	fmt.Fprintf(ctx, "POST xyz=%s <br/>",
		bytes.Join(postValues.PeekMulti("xyz"), []byte(",")))
}

func main() {
	//　使用　fasthttprouter 创建路由
	router := fasthttprouter.New()
	router.GET("/hello/:name", httpHandler)
	log.Println(fasthttp.ListenAndServe(":9090", router.Handler))
}
