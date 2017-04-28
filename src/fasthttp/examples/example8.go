package main

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
	"time"
)

// func httpHandler(ctx *fasthttp.RequestCtx) {
// 	// var req fasthttp.Request
// 	// ctx.Request.CopyTo(&req)
// 	go func() {
// 		time.Sleep(5 * time.Second)
// 		fmt.Println("GET abc=" + string(ctx.URI().QueryArgs().Peek("abc")))
// 		log.Println("done")
// 	}()
// 	ctx.WriteString("hello fasthttp") //立即使用ctx返回数据，而在另一个ｇｏｒｏｕｔｉｎｅ中　对复制出的　ｒｅｑ进行一定的操作
// }

func main() {
	//　使用　fasthttprouter 创建路由
	router := fasthttprouter.New()
	router.GET("/", httpHandler)
	log.Println(fasthttp.ListenAndServe(":9090", router.Handler))
}

// use CopyTo, ctx => req
func httpHandler(ctx *fasthttp.RequestCtx) {
	var req fasthttp.Request
	ctx.Request.CopyTo(&req)
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("GET abc=" + string(req.URI().QueryArgs().Peek("abc")))
		log.Println("done")
	}()
	ctx.WriteString("hello fasthttp") //立即使用ctx返回数据，而在另一个ｇｏｒｏｕｔｉｎｅ中　对复制出的　ｒｅｑ进行一定的操作
}
