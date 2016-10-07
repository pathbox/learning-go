package main

import (
	"io"
	"net/http"

	"github.com/gohttp/app"
	"github.com/gohttp/logger"
)

func main() {
  a := app.New()

  a.Use(logger.New())
  a.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    w.Write([]byte("hello"))
    w.Write([]byte(" world"))
    }))

  a.Get("/baidu", func(w http.ResponseWriter, r *http.Request){
    page, _ := http.Get("http://www.baidu.com")
    defer page.Body.Close()
    io.Copy(w, page.Body)
  })

  a.Get("/error", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("boom"))
	}))

  a.Listen(":9090")

}
