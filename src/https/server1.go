package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,
		"Hi, This is an example of https service in golang!")
	log.Println("handler")
}

func main() {
	fmt.Println("Start server")
	http.HandleFunc("/", handler)
	log.Fatalln(http.ListenAndServeTLS("wsecho.com:9099", "./ca_key/server.crt", "./ca_key/server.key", nil))
}
