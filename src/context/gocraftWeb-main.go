package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gocraft/web"
)

type Context struct {
	CSRFToken string
	User      string
}

// Our middleware *and* handlers must be defined as methods on our context struct,
// or accept the type as their first argument. This ties handlers/middlewares to our
// particular application structure/design.

func (c *Context) CSRFMiddleware(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	token, err := GenerateToken(r)
	if err != nil {
		http.Error(w, "No good!", http.StatusInternalServerError)
		return
	}
	c.CSRFToken = token
	next(w, r)
}

func (c *Context) ShowSignupForm(w web.ResponseWriter, r *web.Request) {
	// No need to type assert it: we know the type.
	// We can just use the value directly.
	fmt.Fprintf(w, "Our token is %v", c.CSRFToken)
}

func main() {
	router := web.New(Context{}).Middleware((*Context).CSRFMiddleware).
		Middleware(web.LoggerMiddleware).
		Middleware(web.ShowErrorsMiddleware)
	router.Get("/signup", (*Context).ShowSignupForm)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}

func GenerateToken(r *web.Request) (string, error) {
	return "our-fake-token", nil
}
