package main

import (
	"fmt"
	"log"
	"net/http"
)

type helloHandler struct {}

func (h helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello, you have hit %s\n", r.URL.Path)
}

func main() {
  log.Fatal(http.ListenAndServe(":9090", helloHandler{}))
}
