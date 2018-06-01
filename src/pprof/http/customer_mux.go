package main

import (
	"net/http"
	"net/http/pprof"
)

func main() {
	// 自定义的mux
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	http.ListenAndServe(":6066", mux)
}

// go tool pprof http://localhost:6066/debug/pprof/heap

// _ "net/http/pprof"
// 如果这样导入包，默认用的mux是 http.DefaultServeMux,如果要使用自定义的mux，则需要去手动定义 mux.HandleFunc("/debug/pprof/", pprof.Index) 等这些HandleFunc
