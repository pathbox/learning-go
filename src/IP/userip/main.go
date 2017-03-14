package main

import (
	"fmt"
	"net"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":9090", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	fmt.Println(query)

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		panic(err)
	}

	userIP := net.ParseIP(ip)
	fmt.Println(userIP)
	if userIP == nil {
		panic(err)
	}
	w.Write([]byte(userIP))
}
