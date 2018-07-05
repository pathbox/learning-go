package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	fmt.Println(":9090")
	http.ListenAndServe(":9090", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("=========", r.Header.Get("Content-Type"))
	r.ParseForm()
	// w.Header().Set("Content-Type", "application/json")
	res := r.PostForm
	fmt.Println("res:", res)
	resp := res.Get("type")
	fmt.Println("resp:", resp)
	w.Write([]byte(resp))
}
