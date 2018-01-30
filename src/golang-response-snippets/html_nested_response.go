package main

import (
	"html/template"
	"path"

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

	fp := path.Join("templates", "index1.html")
	lp := path.Join("templates", "layout1.html")
	tmpl, _ := template.ParseFiles(lp, fp)
	if err := tmpl.Execute(ctx, profile); err != nil {
		ctx.Error("json marshal wrong", fasthttp.StatusInternalServerError)
	}
}
