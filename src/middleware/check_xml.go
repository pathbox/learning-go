package main

import (
	"bytes"
     	"net/http"
)

func enforceXMLHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength == 0 {
			http.Error(w, http.StatusText(400), 400)
			return
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(.Body)
		if http.DetectContentType(buf.Bytes()) != "text/xml; charset=utf-8" {
			http.Error(w, http.StatusText(415), 415)
			return
		}
		next.ServeHTTP(w, r)
		})
}

func main() {
	finalHandler := http.HandlerFunc(final)

	http.Handle("/", enforceXMLHandler(finalHandler))
	http.ListenAndServe(":9090", nil)
}

func final(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
