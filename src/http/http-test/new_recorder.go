package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Something failed", http.StatusInternalServerError)
		// fmt.Println(w)
	}

	req, err := http.NewRequest("GET", "http://httpbin.org/get", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler(w, req)

	fmt.Printf("%d - %s", w.Code, w.Body.String())
}
