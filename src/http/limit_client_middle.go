package main

import (
	"io"
	"log"
	"net/http"
)

func getExpensiveResource() string {
  return "expensive string"
}

// 中间件定义
func maxClients(h http.Handler, n int) http.Handler {
  sema := make(chan struct{}, n)
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    sema <- struct{}{}
    defer func() { <-sema }()
    h.ServeHTTP(w, r)
  })
}

func main() {
  handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    res := getExpensiveResource()
    io.WriteString(w, res)
  })

  http.Handle("/", maxClients(handler, 10)) // http.Handle(url, http.Handler)
  log.Fatal(http.ListenAndServe(":9091", nil))
}
