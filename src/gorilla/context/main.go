package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/context"
	"log"
	"net/http"
)

type contextKey int

// Define keys that support equality.
const csrfKey contextKey = 0
const userKey contextKey = 1

var ErrCSRFTokenNotPresent = errors.New("CSRF token not present in the request context.")

// We'll need a helper function like this for every key:type
// combination we store in our context map else we repeat this
// in every middleware/handler that needs to access the value.
func GetCSRFToken(r *http.Request) (string, error) {
	val, ok := context.GetOk(r, csrfKey)
	if !ok {
		return "", ErrCSRFTokenNotPresent
	}

	token, ok := val.(string)
	if !ok {
		return "", ErrCSRFTokenNotPresent
	}

	return token, nil
}

// A bare-bones example
func CSRFMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// A minimal example: in reality, our CSRF middleware would
		// do a *lot* more than just this.
		maskedToken, err := GenerateToken(r)
		if err != nil {
			http.Error(w, "No good!", http.StatusInternalServerError)
			return
		}

		// The map is global, so we just call the Set function
		context.Set(r, csrfKey, maskedToken)

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func ShowSignupForm(w http.ResponseWriter, r *http.Request) {
	// We'll use our helper function so we don't have to type assert
	// the result in every handler that triggers/handles a POST request.
	csrfToken, err := GetCSRFToken(r)
	if err != nil {
		http.Error(w, "No good!", http.StatusInternalServerError)
		return
	}

	// We can access this token in every handler we wrap with our
	// middleware. No need to set/get from a session multiple times per
	// request (which is slow!)
	fmt.Fprintf(w, "Our token is %v", csrfToken)
}

func main() {
	r := http.NewServeMux()
	r.Handle("/signup", CSRFMiddleware(http.HandlerFunc(ShowSignupForm)))
	// Critical that we call context.ClearHandler here, else we leave
	// old requests in the map.
	err := http.ListenAndServe(":8000", context.ClearHandler(r))
	if err != nil {
		log.Fatal(err)
	}
}

func GenerateToken(r *http.Request) (string, error) {
	return "our-fake-token", nil
}
