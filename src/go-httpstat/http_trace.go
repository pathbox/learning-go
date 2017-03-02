package main

import (
	"github.com/tcnksm/go-httpstat"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	req, err := http.NewRequest("GET", "https://www.udesk.cn", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create a http stat powered context
	var result httpstat.Result
	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	// Send request by deault HTTP client
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
		log.Fatal(err)
	}

	res.Body.Close()
	end := time.Now()
	result.End(end)

	// Show the results
	log.Printf("DNS lookup: %d ms",
		int(result.DNSLookup/time.Millisecond))
	log.Printf("TCP connection: %d ms",
		int(result.TCPConnection/time.Millisecond))
	log.Printf("TLS handshake: %d ms", int(result.TLSHandshake/time.Millisecond))
	log.Printf("Server processing: %d ms", int(result.ServerProcessing/time.Millisecond))
	log.Printf("Content transfer: %d ms", int(result.ContentTransfer(time.Now())/time.Millisecond))
}
