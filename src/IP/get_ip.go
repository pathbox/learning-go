package main

import "net/http"
import (
	"fmt"
	"github.com/sebest/xff"
	"log"
)

func main() {
	log.Println("here start test")
	http.HandleFunc("/get_ip", myHandler)
	log.Fatal(http.ListenAndServe("0.0.0.0:6003", nil))
}

func myHandler(w http.ResponseWriter, req *http.Request) {
	forward := req.Header.Get("X-Forwarded-For")
	log.Println(req.Header)
	fmt.Println(forward)
	fmt.Sprintf("=======%v", forward)
	fmt.Println("======", req.RemoteAddr)
	log.Printf("Get IP from Header: [%v] [%v] [%v] \n [%v]",
		xff.GetRemoteAddr(req), forward, req.RemoteAddr, req.Header)

	w.Write([]byte(req.RemoteAddr))
}
