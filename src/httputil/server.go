package main

import (
	"fmt"
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

	// b, err := ioutil.ReadAll(r.Body)

	// if err != nil {
	// 	w.Write([]byte(string(err.Error())))
	// 	return
	// }

	// fmt.Printf("===request body: %s\n", b)

	// for name, headers := range r.Header {
	// 	fmt.Printf("-name:%v---headers:%v-\n", name, headers)

	// 	for _, h := range headers {
	// 		fmt.Printf("%v:%v\n", name, h)
	// 	}
	// }

	// base64decoder := base64.NewDecoder(base64.StdEncoding, r.Body)
	// gz, err := zlib.NewReader(base64decoder) // 对r.Body 进行base64 编码后压缩

	// if err != nil {
	// 	fmt.Println("===err")
	// 	w.Write([]byte(string(err.Error())))
	// 	return
	// }

	// defer gz.Close()

	// decoder := json.NewDecoder(gz)
	// var t map[string]interface{}

	// err = decoder.Decode(&t)

	// if err != nil {
	// 	w.Write([]byte(string(err.Error())))
	// 	return
	// }

	fmt.Printf("+++req result map: %v\n", t)
	w.Write([]byte("Hello World!"))
}
