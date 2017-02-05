package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %q", html.EscapeString(r.URL.Path))
	})
	log.Fatal(http.ListenAndServe(":9090", nil))
}

/*
这种创建web server的方式，我们称他为“functional”式或者“package level”式的，我们使用了包级别的函数访问隐式的的全局 http 服务和一个隐式的的全局logger实例，我们只用内联函数作为我们handler。

*/
