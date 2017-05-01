package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
	"time"
)

func httpHandler(ctx *fasthttp.RequestCtx) {
	resCh := make(chan string, 1)
	go func() {
		// time.Sleep(5 * time.Second)
		resCh <- string(ctx.FormValue("abc"))
	}()

	select {
	case <-time.After(1 * time.Second):
		ctx.TimeoutError("Timeout")
	case r := <-resCh:
		ctx.WriteString("get: abc = " + r)
	}
	// ctx 在不同的　goroutinue 中
}

func main() {
	//　使用　fasthttprouter 创建路由
	router := fasthttprouter.New()
	router.GET("/", httpHandler)
	log.Println(fasthttp.ListenAndServe(":9090", router.Handler))
}

// RequestCtx 复用引发数据竞争

// RequestCtx 在 fasthttp 中使用 sync.Pool 复用。在执行完了 RequestHandler 后当前使用的 RequestCtx 就返回池中等下次使用。如果你的业务逻辑有跨 goroutine 使用 RequestCtx，那可能遇到：同一个 RequestCtx 在 RequestHandler 结束时放回池中，立刻被另一次连接使用；业务 goroutine 还在使用这个 RequestCtx，读取的数据发生变化。

// 为了解决这种情况，一种方式是给这次请求处理设置 timeout ，保证 RequestCtx 的使用时 RequestHandler 没有结束：
// 还提供 fasthttp.TimeoutHandler 帮助封装这类操作。

// 另一个角度，fasthttp 不推荐复制 RequestCtx。但是根据业务思考，如果只是收到请求数据立即返回，后续处理数据的情况，复制 RequestCtx.Request 是可以的，因此也可以使用

// func httpHandle(ctx *fasthttp.RequestCtx) {
//     var req fasthttp.Request
//     ctx.Request.CopyTo(&req)
//     go func() {
//         time.Sleep(5 * time.Second)
//         fmt.Println("GET abc=" + string(req.URI().QueryArgs().Peek("abc")))
//     }()
//     ctx.WriteString("hello fasthttp")  //立即使用ctx返回数据，而在另一个ｇｏｒｏｕｔｉｎｅ中　对复制出的　ｒｅｑ进行一定的操作
// }
