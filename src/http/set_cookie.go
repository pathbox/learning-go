package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now()
	expiration = expiration.AddDate(0, 0, 1)
	cookie := http.Cookie{Name: "username", Value: "hello world", Expires: expiration}
	http.SetCookie(w, &cookie)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Cookies())
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/get", getHandler)
	log.Fatal(http.ListenAndServe(":9090", nil))
}
