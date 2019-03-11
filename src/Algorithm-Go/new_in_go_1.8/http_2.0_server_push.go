package main

import (
	"net/http"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/",
	http.FileServer(http.Dir("./static"))))
	http.Handle("/index.html", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if p, ok := w.(http.Pusher); ok {
			p.Push("/static/avatar.jpg", nil)
		}
		// load the main page
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<img src="/static/avatar.jpg>`))
	}))
	http.ListenAndServeTLS(":9091", "cert.pem", "key.pem", nil)
}

