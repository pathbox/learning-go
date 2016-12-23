package main

import (
	"github.com/you/somepkg/handler"
	"net/http"
)

func main() {
	db, err := sql.Open("connectionstringhere")
	if err != nil {
		log.Fatal(err)
	}

	// Initialise our app-wide environment with the services/info we need.
	env := &handler.Env{
		DB:   db,
		Port: os.Getenv("PORT"),
		Host: os.Getenv("HOST"),
		// We might also have a custom log.Logger, our
		// template instance, and a config struct as fields
		// in our Env struct.
	}

	http.Handle("/", handler.Handler{env, handler.GetIndex})
	log.Fatal(http.ListenAndServe(":9090", nil))
}
