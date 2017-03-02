package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
)

func main() {
	req, _ := http.NewRequest("GET", "https://baidu.com", nil)
	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("Got Conn: %+v\n", connInfo)
		},
		DNSDone: func(dnsinfo httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Info; %+v\n", dnsinfo)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	_, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Fatal(err)
	}
}
