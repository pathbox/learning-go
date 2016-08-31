package main

import (
  "net/http"
	"github.com/codemodus/catena"
)

func firstWapper(hr http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    _, _ = w.Write([]byte("1"))
    hr.ServeHTTP(w, r)
    _, _ = w.Write([]byte("1"))
  })
}

func endHandler(w http.ResponseWriter, r *http.Request) {
  _, _ = w.Write([]byte("end"))
  return
}

func main() {
  cat0 := catena.New(firstWapper)
  sm := http.NewServeMux()
  sm.Handle("/hello", cat0.EndFn(endHandler))
  sm.ServeHTTP
}
