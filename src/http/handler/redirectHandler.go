package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	rh := http.RedirectHandler("http://www.hao123.com", http.StatusTemporaryRedirect)
	rh_baidu := http.RedirectHandler("http://www.baidu.com", http.StatusTemporaryRedirect)
	mux.Handle("/foo", rh)
	mux.Handle("/baidu", rh_baidu)

	log.Println("Listening at :9090")
	// http.ListenAndServe(":9090", mux)
	log.Fatal(http.ListenAndServe(":9090", mux))
}
