package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/http2"
)

const url = "https://localhost:8080"

var httpVersion = flag.Int("version", 2, "HTTP version")

func main() {
	flag.Parse()
	client := &client.Client{}
	caCert, err := ioutil.ReadFile("server.crt")
	if err != nil {
		log.Fatalf("Reading server certificate: %s", err)
	}
	caCertPool := x509.NewCertPool() // caCertPool
	caCertPool.AppendCertsFromPEM(caCert)
	// Create TLS configuration with the certificate of the server
	tlsConfig := &tls.Config{
		RootCAs: caCertPool, // caCertPool
	}

	// Use the proper transport in the client
	switch *httpVersion {
	case 1:
		client.Transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	case 2:
		client.Transport = &http2.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	resp, err := client.Get(url)
	defer resp.Body.Close()

	body, err :=ioutil.ReadAll(resp.Body)
	if err != nil {
 		log.Fatalf("Failed reading response body: %s", err)
	}
	fmt.Printf(
		"Got response %d: %s %s\n",
		resp.StatusCode, resp.Proto, string(body)
        )
}
