package main

import (
	"github.com/fuxiaohei/GoInk"
)

func main() {
	app := GoInk.New()
	app.Get("/", func(ctx *GoInk.Context) {
		ctx.Body = []byte("Hello World!")
	})
	app.Run()
}
