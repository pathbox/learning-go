package main

import (
	"encoding/json"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type Profile struct {
	Name    string
	Hobbies []string
}

func main() {
	router := fasthttprouter.New()
	router.GET("/", Index)
	fasthttp.ListenAndServe(":9000", router.Handler)
}

func Index(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Color", "green green")
	profile := Profile{"Alex", []string{"snowboarding", "programming"}}

	js, err := json.Marshal(profile)

	if err != nil {
		ctx.Error("json marshal wrong", fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetContentType("application/json")
	ctx.Write(js)
}
