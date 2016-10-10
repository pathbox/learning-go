package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/", Index)
  log.Fatal(http.ListenAndServe(":9090", router))
}

func Index(w http.ResponseWriter, r *http.Request)  {
  fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
