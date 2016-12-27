package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func AddContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, "-", r.RequestURI)
		cookie, _ := r.Cookie("username")
		if cookie != nil {
			ctx := context.WithValue(r.Context(), "Username", cookie.Value)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func StatusPage(w http.ResponseWriter, r *http.Request) {
	if username := r.Context().Value("Username"); username != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello " + username.(string) + "\n"))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Logged in"))
	}
}

func LogininPage(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().Add(365 * 24 * time.Hour) // Set to expire in 1 year
	cookie := http.Cookie{Name: "username", Value: "alice_cooper@gmail.com", Expires: expiration}
	http.SetCookie(w, &cookie)
}

func LoginoutPage(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().AddDate(0, 0, -1) // Set to expire in the past
	cookie := http.Cookie{Name: "username", Value: "alice_cooper@gmail.com", Expires: expiration}
	http.SetCookie(w, &cookie)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", StatusPage)
	mux.HandleFunc("/login", LogininPage)
	mux.HandleFunc("/logout", LoginoutPage)

	log.Println("Start server on port: 9090")
	contextedMux := AddContext(mux)

	log.Fatal(http.ListenAndServe(":9090", contextedMux))
}
