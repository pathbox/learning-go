package main

import "net/http"

func main() {
	// http://localhost:8888/static
	// http.Handle("/static", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	fileHandler := http.FileServer(http.Dir("static"))
	http.Handle("/", fileHandler)
	http.ListenAndServe(":9090", nil)
}
