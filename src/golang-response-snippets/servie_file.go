package main

import (
	"path"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func main() {
	router := fasthttprouter.New()
	router.GET("/", Index)
	fasthttp.ListenAndServe(":9000", router.Handler)
}

func Index(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Color", "green green")

	fp := path.Join("images", "seabar.jpg")

	fasthttp.ServeFile(ctx, fp)
}
