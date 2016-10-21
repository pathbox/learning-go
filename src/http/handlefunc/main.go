package main

import (
  // "fmt"
  "net/http"
)

func main() {
  myHandle := func(w http.ResponseWriter, r *http.Request) {
    // fmt.Fprintf(w, "Hello World")
    w.Write([]byte("Hello World"))
  }

  http.HandleFunc("/", myHandle)
  http.ListenAndServe(":9090", nil)
}
