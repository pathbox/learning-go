package main

import (
	"log"
	"net/http"
	"time"
)

type timeHandler struct {
	zone *time.Location
}

func (th *timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().In(th.zone).Format(time.RFC1123)
	w.Write([]byte("The time is: " + tm))
}

func newTimeHandler(name string) *timeHandler {
	return &timeHandler{zone: time.FixedZone(name, 8*3600)}
}

func main() {
	myHandler := newTimeHandler("CST")
	//Custom http server
	server := &http.Server{
		Addr:           ":8080",
		Handler:        myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(server.ListenAndServe())
}
