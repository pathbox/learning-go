package main
import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,
        "Hi, This is an example of https service in golang!")
}

func main() {
	fmt.Println("Start server")
	http.HandleFunc("/", handler)
	http.ListenAndServeTLS("127.0.0.1:9099", "./ca_key/server.crt", "./ca_key/server.key", nil)
}