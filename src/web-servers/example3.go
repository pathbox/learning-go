package main

import (
	"io"
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello world\n")
}

func main() {
	http.HandleFunc("/hello", HelloServer)
	log.Fatal(http.ListenAndServe(":9090", nil))
}
