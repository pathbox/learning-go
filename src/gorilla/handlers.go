package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func main() {
	finalHandler := http.HandlerFunc(final)

	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	http.Handle("/", handlers.LoggingHandler(logFile, finalHandler))
	http.ListenAndServe(":9000", nil)
}

func final(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
