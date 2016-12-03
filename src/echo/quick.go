package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"net/http"
)

func main() {
	e := echo.New()
//	e.Use(middleware.Logger())
	e.Get("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})
	e.Run(standard.New(":9090"))
}
