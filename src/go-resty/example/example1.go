package main

import (
	"gopkg.in/resty.v0"
	"log"
)

func main() {

	resp, err := resty.R().Get("http://httpbin.org/get")
	log.Printf("\nError: %v", err)
	log.Printf("\nResponse Status Code: %v", resp.StatusCode())
	log.Printf("\nResponse Status: %v", resp.Status())
	log.Printf("\nResponse Time: %v", resp.Time())
	log.Printf("\nResponse Recevied At: %v", resp.ReceivedAt())
	log.Printf("\nResponse Body: %v", resp) // or resp.String() or string(resp.Body())

}
