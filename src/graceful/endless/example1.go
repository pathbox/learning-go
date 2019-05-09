package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fvbock/endless"
	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	duration, err := time.ParseDuration(r.FormValue("duration"))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	time.Sleep(duration)

	pid := os.Getpid()

	re := fmt.Sprintf("New Version:%d", pid)
	w.Write([]byte(re))
}

func main() {
	mux1 := mux.NewRouter()
	mux1.HandleFunc("/hello", handler).
		Methods("GET")

	err := endless.ListenAndServe("localhost:4242", mux1)
	if err != nil {
		log.Println(err)
	}
	log.Println("Server on 4242 stopped")

	os.Exit(0)
}
