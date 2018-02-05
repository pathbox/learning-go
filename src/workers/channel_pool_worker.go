package main

import (
	"fmt"
	"math/rand"
	"runtime"

	_ "net/http/pprof"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

var Pool chan int32

var C int

func main() {
	runtime.GOMAXPROCS(4)
	Pool = make(chan int32, 10000)
	go doJob()
	router := fasthttprouter.New()
	router.GET("/pool", WorkerAction)
	fasthttp.ListenAndServe(":9090", router.Handler)
}

func WorkerAction(ctx *fasthttp.RequestCtx) {
	r := rand.Int31n(100000)
	Pool <- r
	ctx.Write([]byte("OK"))
}

func doJob() {
	for {
		select {
		case <-Pool:
			// fmt.Println(result)
			C++
			fmt.Println("Count: ", C)
		}
	}
}

//开启一个goroutine，作为 doJob 处理。循环从Pool中取出数据（数据以后可以换成定义的struct），进行处理
// 类似 简单的事件驱动的模式， 这样异步处理中，只会有一个异步线程在处理
// 与multi的方式不同，multi的方式会产生大量的 goroutine进行异步处理
// 但这里 使用了 全局的 Pool channel，作为 队列。会消耗一些内存，和 响应延迟
// 但通过 ab测试， 发现比multi方式每秒处理请求数量多1000
