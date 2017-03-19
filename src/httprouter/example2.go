package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func BasicAuth(h httprouter.Handle, requireUser, requirePassword string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, password, hasAuth := r.BasicAuth()

		if hasAuth && user == requireUser && password == requirePassword {
			h(w, r, ps)
			log.Println("nice nice")
		} else {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("Not Protected\n"))
}

func Protected(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("Protected\n"))
}

func main() {
	user := "admin"
	password := "password"

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/protected/", BasicAuth(Protected, user, password))

	log.Fatal(http.ListenAndServe(":9090", router))
}
