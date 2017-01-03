package main

import (
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func main() {
	http.HandleFunc("/", helloHandler)
	http.ListenAndServe(":9090", nil)
}

// 其实，HandleFunc 只是一个适配器，

// // The HandlerFunc type is an adapter to allow the use of
// // ordinary functions as HTTP handlers.  If f is a function
// // with the appropriate signature, HandlerFunc(f) is a
// // Handler object that calls f.
// type HandlerFunc func(ResponseWriter, *Request)

// // ServeHTTP calls f(w, r).
// func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
//   f(w, r)
// }
// 自动给 f 函数添加了 HandlerFunc 这个壳，最终调用的还是 ServerHTTP，只不过会直接使用 f(w, r)。这样封装的好处是：使用者可以专注于业务逻辑的编写，省去了很多重复的代码处理逻辑。如果只是简单的 Handler，会直接使用函数；如果是需要传递更多信息或者有复杂的操作，会使用上部分的方法
