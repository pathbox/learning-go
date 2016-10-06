package example

import (
	"io"
	"net/http"
)

func basicHelloHandler(w http.ResponseWriter, r *http.Request)  {
  io.WriteString(w, "Hello World")
}

func basicHelloHandler()  {
  http.HandleFunc("/hello", basicHelloHandler)
}

func BasicEngine() http.Handler  {
  mux := http.NewServeMux()
  mux.HandleFunc("/", basicHelloHandler)

  return mux
}
