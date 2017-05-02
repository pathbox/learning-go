package main

import (
	"fmt"
	"gopkg.in/resty.v0"
	"os"
	"time"
)

func main() {
	// Using your custom log writer
	logFile, _ := os.OpenFile("./restry.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	resty.SetLogger(logFile)

	resty.SetTimeout(time.Duration(1 * time.Minute))

	resp, err := resty.R().Get("http://httpbin.org/get")
	fmt.Printf("\nError: %v", err)
	fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())
	fmt.Printf("\nResponse Status: %v", resp.Status())
	fmt.Printf("\nResponse Time: %v", resp.Time())
	fmt.Printf("\nResponse Recevied At: %v", resp.ReceivedAt())
	fmt.Printf("\nResponse Body: %v", resp) // or resp.String() or string(resp.Body())

}
