package main

import (
	"./handler"
	"log"
	"net/http"
  "os"
)

func main() {
	db, err := sql.Open("connectionsstringhere")

	if err != nil {
		log.Fatal(err)
	}

  env := &handler.Env{
    DB: db,
    Port: os.Getenv("PORT"),
    Host: os.Getenv("HOST")
  }

  http.Handle("/", handler.Handler{env, handler.GerIndex})
  log.Fatal(http.ListenAndServe(":9090", nil))
}
