package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/http2"
)

func main() {
	url := "http://localhost:9090/"
	client(url)
}

func client(url string) {
	log.SetFlags(log.Llongfile)
	tr := &http2.Transport{
		AllowHTTP: true, //充许非加密的链接
		// TLSClientConfig: &tls.Config{
		//     InsecureSkipVerify: true,
		// },
		DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(netw, addr)
		},
	}

	httpClient := http.Client{Transport: tr}

	resp, err := httpClient.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("resp StatusCode:", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("resp.Body:\n", string(body))
}
