package main

import (
	"net/http"

	"github.com/gohttp/app"
	"github.com/gohttp/logger"
	"github.com/gohttp/mount"
	"github.com/gohttp/serve"
)

func main() {
  a := app.New()
  a.Use(logger.New())
  a.Use(mount.New("/example", serve.New("example")))
  a.Use(mount.New("/blog", blog()))
  a.Use(mount.New("/hello", hello))
  a.Get("/", http.HandlerFunc(hello))

  a.Listen(":9090")
}

func blog() *app.App{
  a := app.New()
  a.Get("", func(w http.ResponseWriter, r *http.Request){
    w.Write([]byte("posts\n"))
  })
  return a
}

func hello(w http.ResponseWriter, r *http.Request){
  w.Write([]byte("hello\n"))
}
