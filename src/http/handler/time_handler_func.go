package main

import (
	"log"
	"net/http"
	"time"
)

func timeHandler(format string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tm := time.Now().Format(format)
		w.Write([]byte("The time is: " + tm))
	}
}

func main() {
	th := timeHandler(time.RFC1123)
	log.Println(&th)
	http.Handle("/", th)
	log.Println("Listening at :9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
