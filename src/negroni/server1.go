package main

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/home", HomeHandler)

	// n := negroni.New(middleware1, middleware2...)
	n := negroni.Classic()
	// n.Use(Middleware3)
	n.UseHandler(router)
	n.Run(":9090")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome Home!"))
}
