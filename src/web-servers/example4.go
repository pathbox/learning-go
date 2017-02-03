package main

import (
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server. \n"))
}

func main() {
	http.HandleFunc("/", handler)
	log.Printf("About to listen on 9090. Go to https://127.0.0.1:9090")
	// err := http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
	err := http.ListenAndServe(":9090", nil)
	log.Fatal(err)
}
