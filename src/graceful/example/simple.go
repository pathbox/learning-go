package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"gopkg.in/tylerb/graceful.v1"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := os.Getpid()
		fmt.Fprintf(w, "Welcome to the home page! Pid: %d", p)
	})
	graceful.Run(":9091", 10*time.Second, mux)
}
