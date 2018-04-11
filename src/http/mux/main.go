package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars)
	fmt.Println("++++++", vars["id"])

	fmt.Fprintln(w, "hello")
}

func main() {
	handlers := mux.NewRouter()
	handlers.HandleFunc("/names/1", handler)

	handlers.HandleFunc("/2", handler)
	handlers.HandleFunc("/agents/{id:[0-9]+}", handler)
	server := &http.Server{
		Addr:         ":8080",
		Handler:      handlers,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	server.ListenAndServe()
}
