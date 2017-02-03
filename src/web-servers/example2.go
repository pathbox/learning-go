package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		log.Println("request to foo")
	})
	log.Fatal(http.ListenAndServe(":9090", nil))
}
