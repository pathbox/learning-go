package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func main() {
	finalHandler := http.HandlerFunc(final) // 定义一个handler

	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	http.Handle("/", handlers.LoggingHandler(logFile, finalHandler)) // handler的嵌套
	http.ListenAndServe(":9090", nil)
}

func final(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
