package main

import (
	"log"
	"net/http"
	"time"
)

func timeHandler(w http.ResponseWriter, r *http.Request) { // 实现w 和 r 的方法
	tm := time.Now().Format(time.RFC1123)
	w.Write([]byte("The time is: " + tm))
}

func main() {

	http.HandleFunc("/time", timeHandler) // 直接使用http.HandleFunc，其底层会用 HandlerFunc 对timeHandler包装后返回

	log.Println("Listening...")

	log.Fatal(http.ListenAndServe(":9090", nil)) // 使用的是 DefaultServerMux
}
