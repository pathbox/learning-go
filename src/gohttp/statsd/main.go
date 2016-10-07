package main

import (
	"net/http"
	"time"

	"github.com/gohttp/app"
	stats "github.com/gohttp/statsd"
)

func main() {
	a := app.New()

	statsd, _ := statsd.Dial(":8125")

	a.Use(stats.New(statsd))

	a.Get("/", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("hello"))
		res.Write([]byte(" world"))
	}))

	a.Get("/slow", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		time.Sleep(time.Second)
		res.Write([]byte("hello"))
		res.Write([]byte(" world"))
	}))

	a.Listen(":3000")
}
