package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
  expiration := time.Now()
  expiration = expiration.AddDate(1, 0, 0)
  cookie := http.Cookie{Name: "username", Value: "Akagi201", Expires: expiration}
  http.SetCookie(w, &cookie)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println(r.Cookies())
}

func main() {
  http.HandleFunc("/", rootHandler)
  http.HandleFunc("/getcookie", getHandler)

  log.Fatal(http.ListenAndServe(":9090", nil))
}
