package handler

import (
	"net/http"
)

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
}

//StatusError represents an error with an associated HTTP status code.

type StatusError struct {
	Code int
	Err  error
}

// Allows StatusError to satisfy the error interface.
func (se StatusError) Status() int {
	return se.Code
}

// A (simple) example of our application-wide configuration.
type Env struct {
	DB   *sql.DB
	Port string
	Host string
}

// The Handler struct that takes a configured Env and a function matching
// our useful signature.

type Handler struct {
	*Env
	H func(e *Env, w http.ResponseWriter, r *http.Request) error
}

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.Env, w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

func GetIndex(env *Env, w http.ResponseWriter, r *http.Request) error {
	users, err := env.DB.GetAllUsers()
	if err != nil {
		// We return a status error here, which conveniently wraps the error
		// returned from our DB queries. We can clearly define which errors
		// are worth raising a HTTP 500 over vs. which might just be a HTTP
		// 404// We return a status error here, which conveniently wraps the error
		// returned from our DB queries. We can clearly define which errors
		// are worth raising a HTTP 500 over vs. which might just be a HTTP
		// 404, 403 or 401 (as appropriate). It's also clear where our
		// handler should stop processing by returning early., 403 or 401 (as appropriate). It's also clear where our
		// handler should stop processing by returning early.
		return StatusError{500, err}
	}
	fmt.Fprintf(w, "%+v", users)
	return nil
}
