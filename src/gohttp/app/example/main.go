package main

import (
	"fmt"
	"net/http"
  "io"

	"github.com/gohttp/app"
)

func main() {
  app := app.New()
  app.Get("/foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    io.WriteString(w, "hello world")
    }))

    app.Get("/bar", func(w http.ResponseWriter, r *http.Request){
      fmt.Fprintln(w, "bar")
    })
    
    app.Listen(":9090")
}
