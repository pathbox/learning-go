package main

import (
	"fmt"
	"log"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

func main() {

	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query()

	fmt.Println(p)
	fmt.Println(p["name"])

	b, _ := jsoniter.Marshal(p)

	fmt.Println(string(b))

	w.Write([]byte("OK"))
}
