package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func main() {
	router := fasthttprouter.New()
	router.GET("/", Index)
	fasthttp.ListenAndServe(":9000", router.Handler)
}

func Index(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Server", "A Go Web Server")
	ctx.SetStatusCode(fasthttp.StatusOK)
}
