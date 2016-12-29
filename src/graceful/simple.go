package main

import (
	"fmt"
	"gopkg.in/tylerb/graceful.v1"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})
	graceful.Run(":9090", 10*time.Second, mux)
}
