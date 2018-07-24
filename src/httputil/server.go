package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

func main() {
	http.HandleFunc("/index", index)
	fmt.Println("listening 9090")
	http.ListenAndServe(":9090", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	reqDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		w.Write([]byte(string(err.Error())))
		return
	}

	fmt.Printf("===The request content:%s\n", reqDump)

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.Write([]byte(string(err.Error())))
		return
	}

	fmt.Printf("===request body: %s\n", b)

	w.Write([]byte("Hello World!"))
}
