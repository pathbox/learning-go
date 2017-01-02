package main

import (
	// Standard library packages
	"encoding/json"
	"fmt"
	"net/http"

	// Third party packages
	"github.com/julienschmidt/httprouter"
)

type (
	// User represents the structure of our resource
	User struct {
		Name   string
		Gender string
		Age    int
		Id     string
	}
)

func main() {
	// Instantiate a new router
	r := httprouter.New()

	// Get a user resource
	r.GET("/user/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// Stub an example user
		u := User{
			Name:   "Bob Smith",
			Gender: "male",
			Age:    50,
			Id:     p.ByName("id"),
		}

		// Marshal provided interface into JSON structure
		uj, _ := json.Marshal(u)

		// Write content-type, statuscode, payload
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", uj)
	})

	// Fire up the server
	http.ListenAndServe("localhost:3000", r)
}
