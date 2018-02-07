// 每个请求，都起一个goroutine 异步处理任务。 这样相当于多线程异步处理
// 优点简单，缺点 当请求量大的时候，会产生巨大量的goroutine，创建大量的goroutine会成消耗大量资源，成为瓶颈
package main

import (
	"fmt"
	"math/rand"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func main() {
	router := fasthttprouter.New()
	router.GET("/multi", WorkerAction)
	fasthttp.ListenAndServe(":9090", router.Handler)
}

func WorkerAction(ctx *fasthttp.RequestCtx) {
	go doJob()
	ctx.Write([]byte("OK"))
}

func doJob() {
	rand := rand.Int31n(100000)
	fmt.Println(rand)
}
