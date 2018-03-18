package main

import (
	"log"
	"net/http"
	"time"
)

func timeHandler(format string) http.Handler { // 实现w 和 r 的方法
	fn := func(w http.ResponseWriter, r *http.Request) {
		tm := time.Now().Format(format)
		w.Write([]byte("The time is: " + tm))
	}
	return http.HandlerFunc(fn)
}

func main() {

	mux := http.NewServeMux() // 显示定义

	th := timeHandler(time.RFC1123)

	mux.Handle("/time", th) // 直接使用http.HandleFunc，其底层会用 HandlerFunc 对timeHandler包装后返回

	log.Println("Listening...")

	log.Fatal(http.ListenAndServe(":9090", mux)) // 使用的是显示定义的http.NewServeMux()
}
