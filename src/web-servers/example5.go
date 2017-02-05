package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

var (
	logger *log.Logger
)

func main() {
	logger = log.New(os.Stdout, "web ", log.LstdFlags)

	server := &http.Server{
		Addr:    ":9090",
		Handler: routes(),
	}
	server.ListenAndServe()
}

func routes() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("/foo", foo)

	return r
}

func foo(w http.ResponseWriter, r *http.Request) {
	logger.Println("request to foo")
	io.WriteString(w, "Hello World\n")
	w.Write([]byte("This is an example server. \n"))
}
