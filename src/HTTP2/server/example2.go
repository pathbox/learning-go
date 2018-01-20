package main

import (
	"fmt"
	"golang.org/x/net/http2"
	"net/http"

	"net"
	"time"
)

//net/http包默认可以采用http2进行服务，在没有进行https的服务上开启H2，
//需要修改ListenAndServer的默认h2服务

type serverHandler struct {
}

func (sh *serverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	w.Header().Set("server", "h2Server")
	w.Write([]byte("This is a http2 server"))
}

func main() {
	server := &http.Server{
		Addr:         ":9090",
		Handler:      &serverHandler{},
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	server2 := &http2.Server{
		IdleTimeout: 1 * time.Minute,
	}
	http2.ConfigureServer(server, server2)

	l, _ := net.Listen("tcp", ":9090")
	defer l.Close()
	// 不使用默认的 ListenAndServe
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Accept error: ", err)
			continue
		}

		go server2.ServeConn(conn, &http2.ServeConnOpts{BaseConfig: server}) // 每来一个conn， 开启一个goroutine对conn进行读写操作
	}
}

// http2: server: error reading preface from client 127.0.0.1:50871: bogus greeting "GET / HTTP/1.1\r\nHost: 12"
// 说明 HTTP2 需要客户端（比如浏览器）和服务端同时支持
