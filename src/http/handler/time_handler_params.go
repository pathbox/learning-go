package main

import (
	"log"
	"net/http"
	"time"
)

//func timeHandler(format string) http.Handler {
//	fn := func(w http.ResponseWriter, r *http.Request) {
//		tm := time.Now().Format(format)
//		w.Write([]byte("The time is: " + tm))
//	}
//
//	return http.HandlerFunc(fn)
//}

//func timeHandler(format string) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		tm := time.Now().Format(format)
//		w.Write([]byte("The time is: " + tm))
//	})
//}

func timeHandler(format string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tm := time.Now().Format(format)
		w.Write([]byte("The time is: " + tm))
	}
}

func main() {
	mux := http.NewServeMux()

	th := timeHandler(time.RFC1123)
	mux.Handle("/time", th)

	log.Println("Listening at :3000")
	http.ListenAndServe(":3000", mux)
}
