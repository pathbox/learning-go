package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func myLoggingHandler(h http.Handler) http.Handler {
	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		panic(err)
	}
	return handlers.LoggingHandler(logFile, h)
}

func main() {
	finalHandler := http.HandlerFucn(final) // 定义一个handler

	http.Handle("/", myLoggingHandler(finalHandler))
	http.ListenAndServe(":9090", nil)
}

func finalHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
