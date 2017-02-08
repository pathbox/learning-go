package main

import (
	"github.com/goji/httpauth"
	"github.com/zenazn/goji"
	"net/http"
)

func main() {
	goji.Use(httpauth.SimpleBasicAuth("admin", "password"))
	goji.Use(SomeOtherMiddleware)
	goji.Get("/thing", myHandler)
	goji.Serve()
}
