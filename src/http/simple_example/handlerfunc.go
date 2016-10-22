package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
  myhandlerFunc := http.HandlerFunc(myHandler)
  log.Fatal(http.ListenAndServe(":9090", myhandlerFunc))
}

func myHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "hello, you have hit %s\n", r.URL.Path)
}
