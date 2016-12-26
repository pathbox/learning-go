package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

var ErrTypeNotPresent = errors.New("Expected type not present in the request context.")

// A little simpler: we just need this for every *type* we store.
func GetContextString(c web.C, key string) (string, error) {
	val, ok := c.Env[key].(string)
	if !ok {
		return "", ErrTypeNotPresent
	}
	return val, nil
}

// A bare-bones example
func CSRFMiddleware(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		maskedToken, err := GenerateToken(r)
		if err != nil {
			http.Error(w, "No good!", http.StatusInternalServerError)
			return
		}

		// Goji only allocates a map when you ask for it.
		if c.Env == nil {
			c.Env = make(map[interface{}]interface{})
		}
		// Not a global - a reference to the context map
		// is passed to our handlers explicitly.
		c.Env["csrf_token"] = maskedToken
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// Goji's web.HandlerFunc type is an extension of net/http's http.HandlerFunc,
// except it also passes in a request context (aka web.C.Env)

func ShowSignupForm(c web.C, w http.ResponseWriter, r *http.Request) {
	// We'll use our helper function so we don't have to type assert
	// the result in every handler.
	GetCSRFToken, err := GetContextString(c, "csrf_token")
	if err != nil {
		http.Error(w, "No good!", http.StatusInternalServerError)
		return
	}

	// We can access this token in every handler we wrap with our
	// middleware. No need to set/get from a session multiple times per
	// request (which is slow!)
	fmt.Fprintf(w, "Our token is %v", GetCSRFToken)
}

func main() {
	goji.Get("/signup", ShowSignupForm)
	goji.Use(CSRFMiddleware)

	goji.Serve()
}

func GenerateToken(r *http.Request) (string, error) {
	return "our-fake-token", nil
}
