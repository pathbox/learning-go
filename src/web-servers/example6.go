package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	log.Fatal(http.ListenAndServe(":9090", New()))
}

func New() http.Handler {
	mux := http.NewServeMux()
	log := log.New(os.Stdout, "web ", log.LstdFlags)
	app := &app{mux, log}
	mux.HandleFunc("/foo", app.foo)

	return app
}

type app struct {
	mux *http.ServeMux
	log *log.Logger
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

/*
当删除 func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) 的时候
./example6.go:19: cannot use app (type *app) as type http.Handler in return argument:
  *app does not implement http.Handler (missing ServeHTTP method)
*/

func (a *app) foo(w http.ResponseWriter, r *http.Request) {
	a.log.Println("request to foo")
}
