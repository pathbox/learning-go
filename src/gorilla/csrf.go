package main

import (
    "net/http"

    // Don't forget to `go get github.com/gorilla/csrf`
    "github.com/gorilla/csrf"
    "github.com/gorilla/mux"
)

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/signup", showSignupForm)
  CSRF := csrf.Protect([]byte("32-byte-long-auth-key"))
  // PS: Don't forget to pass csrf.Secure(false) if you're developing locally
  // over plain HTTP (just don't leave it on in production).
  log.Fatal(http.ListenAndServe(":9090", CSRF(r)))
}

func ShowSignupForm(w http.ResponseWriter, r *http.Request) {
    // signup_form.tmpl just needs a  template tag for
    // csrf.TemplateField to inject the CSRF token into. Easy!
    t.ExecuteTemplate(w, "signup_form.tmpl", map[string]interface{
        csrf.TemplateTag: csrf.TemplateField(r),
    })
}

func SubmitSignupForm(w http.ResponseWriter, r *http.Request) {
    // We can trust that requests making it this far have satisfied
    // our CSRF protection requirements.
}
