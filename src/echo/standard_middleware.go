package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"net/http"
)

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println("standard middleware")
		next.ServeHTTP(w, r)
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	println("standard handler")
	w.Write([]byte("Hello World!"))
}

func main() {
	e := echo.New()
	e.Use(standard.WrapMiddleware(middleware))
	e.GET("/", standard.WrapHandler(http.HandlerFunc(handler)))
	e.Run(standard.New(":9090"))
}
