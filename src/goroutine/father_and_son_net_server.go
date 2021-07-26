package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/father", father)
	fmt.Println("Listen 9090")
	http.ListenAndServe(":9090", nil)
}

func father(w http.ResponseWriter, r *http.Request) {
	go son()
	fmt.Println("This is father")
	w.Write([]byte("Thie is father"))
}

func son() {
	time.Sleep(5 * time.Second)
	fmt.Println("This is son")
}

// This is father
// This is son
// son goroutinue is OK
