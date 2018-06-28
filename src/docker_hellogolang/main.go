package main

import (
	"log"
	"net/http"
)

func main() {
	addr := ":9097"
	http.HandleFunc("/hello", hello)
  log.Println("Listening 9097")
	log.Fatal(http.ListenAndServe(addr, nil))
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("Hello Golang")
	w.Write([]byte(`Hello Golang`))
}

// CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
