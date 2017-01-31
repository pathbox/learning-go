package main

import (
	"github.com/hoisie/web"
)

func main() {
	web.Get("/", func() string {
		return "Hello Golang"
	})
	web.Run(":8080")
}
