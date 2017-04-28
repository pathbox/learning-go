package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
)

// RequestHandler 类型， 使用RequestCtx 传递HTTP 的数据
func httpHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hello fasthttp") // *RequestCtx 实现了 io.Writer
}

func main() {
	// 一定要写httpHandler 否则会有nil pointer 的错误，没有处理http数据的函数
	log.Println(fasthttp.ListenAndServe(":9090", httpHandler))
}
