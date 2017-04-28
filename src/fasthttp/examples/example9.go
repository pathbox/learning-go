package main

import (
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func httpHandler(ctx *fasthttp.RequestCtx) {
	b := fasthttp.AcquireByteBuffer()
	b.B = append(b.B, "Hello "...)

	b.B = fasthttp.AppendHTMLEscape(b.B, "<strong>world</strong>") // 把　html标签编码，防止攻击
	defer fasthttp.ReleaseByteBuffer(b)                            // 最后记得释放

	ctx.Write(b.B)
}

func main() {
	//　使用　fasthttprouter 创建路由
	router := fasthttprouter.New()
	router.GET("/", httpHandler)
	log.Println(fasthttp.ListenAndServe(":9090", router.Handler))
}
