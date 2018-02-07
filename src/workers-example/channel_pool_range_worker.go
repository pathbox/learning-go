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
	go doRangeJob() // 开启一个守护异步 goroutine
	router := fasthttprouter.New()
	router.GET("/pool", WorkerAction)
	fasthttp.ListenAndServe(":9090", router.Handler)
}

func WorkerAction(ctx *fasthttp.RequestCtx) {
	r := rand.Int31n(100000)
	Pool <- r
	ctx.Write([]byte("OK"))
}

func doRangeJob() {

	for range Pool {
		C++
		fmt.Println("Count: ", C)
	}
}

// 使用range Pool 方式而不是select的方式从Pool中得到数据
// 这种方式适合简单的处理，缺少了 select的调度功能的支持
