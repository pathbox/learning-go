package main

import (
  "fmt"
  "log"
  "net/http"
  "runtime"
)

func main() {
  runtime.GOMAXPROCS( runtime.NumCPU() - 1)

  http.HandleFunc("/", addHeader)
  log.Println("Go http Server Listening at :9090")
  log.Fatal(http.ListenAndServe(":9090", nil))
}

func addHeader(w http.ResponseWriter, r *http.Request) {
  w.Header().Add("Last-Modified", "Thu, 18 Jun 2015 10:10:10 GMT")
  w.Header().Add("Accept-Ranges", "bytes")
  w.Header().Add("E-Tag", "55829c5b-17")
  w.Header().Add("Server", "golang-http-server")
  | fmt.Fprint(w, "<h1>\n Hello World! \n</h1>\n")
}
