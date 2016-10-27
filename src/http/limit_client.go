package main

import (
	"io"
	"log"
	"net/http"
)

func getExpensiveResource() string {
  return "expensive string"
}

func main() {
  const maxClients = 10
  sema := make(chan struct{}, maxClients)

  http.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    sema <- struct{}{}
    defer func() { <- sema }()
    res := getExpensiveResource()
    io.WriteString(w, res)
    }))
    log.Fatal(http.ListenAndServe(":9091", nil))
}
