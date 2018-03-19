package main

import (
	"log"
	"net/http"
	"time"
)

type timeHandler struct {
	format string
}

func (th *timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) { // 实现 ServeHTTP接口
	tm := time.Now().Format(th.format)
	w.Write([]byte("The time is: " + tm))
}

func main() {

	mux := http.NewServeMux() // 显示定义

	th := &timeHandler{format: time.RFC1123}

	mux.Handle("/time", th) // 直接使用http.HandleFunc，其底层会用 HandlerFunc 对timeHandler包装后返回

	log.Println("Listening...")

	log.Fatal(http.ListenAndServe(":9090", mux)) // 使用的是显示定义的http.NewServeMux()
}
