package main

import (
  "fmt"
	"io"
	"net/http"
	"os"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
  io.WriteString(w, "hello world")
}

func BasicEngine() http.Handler {
  mux := http.NewServeMux()
  mux.HandleFunc("/hello", helloHandler)
  return mux
}

func main() {
  mux := BasicEngine()
  if err := http.ListenAndServe(":9090", mux); err != nil{
    fmt.Println(nil, "http listen at 2016 failed. err is", err)
		os.Exit(-1)
  }
}
