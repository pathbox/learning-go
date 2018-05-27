package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", Index)
	fmt.Println("server 63450")
	http.ListenAndServe(":63450", nil)

}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here is 63450 server")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Hello World! Here is 63450"))
}
