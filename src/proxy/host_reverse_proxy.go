package main

import (
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewMultipleHostReverseProxy(targets []*url.URL) *httputil.ReverseProxy {

	director := func(req *http.Request) {
		i := rand.Int()
		log.Printf("rand i:%d\n", i)
		target := targets[i%len(targets)]
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
	}
	return &httputil.ReverseProxy{Director: director} // ReverseProxy help get one http server from director
}

func main() {
	proxy := NewMultipleHostReverseProxy([]*url.URL{
		{
			Scheme: "http",
			Host:   "localhost:9091",
		},
		{
			Scheme: "http",
			Host:   "localhost:9092",
		},
	})
	go s1()
	go s2()
	log.Fatal(http.ListenAndServe(":9090", proxy))
}

func s1() {
	m := http.NewServeMux()
	m.HandleFunc("/", randRes1)
	log.Println("s1")
	log.Fatal(http.ListenAndServe(":9091", m))

}

func s2() {
	m := http.NewServeMux()
	m.HandleFunc("/", randRes2)
	log.Println("s2")
	log.Fatal(http.ListenAndServe(":9092", m))

}

func randRes1(w http.ResponseWriter, r *http.Request) {
	rs := "here is s1"
	w.Write([]byte(rs))
}

func randRes2(w http.ResponseWriter, r *http.Request) {
	rs := "here is s2"
	w.Write([]byte(rs))
}
