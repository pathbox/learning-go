package main

import(
	"flag"
	"github.com/levigross/grequests"
	"encoding/json"
	"fmt"
)

var count int

func main() {
	flag.IntVar(&count,"count", 10, "goroutinue count")

	for i :=0; i <= count; i++{
		go doRequest()
	}

  select{}
}

func doRequest() {
	url := "127.0.0.1:3344/go/push"
	j,_ := json.Marshal([]byte(`{"TextMsg": "hello morning2", "DeviceToken": "56a1921b464128c031750e2eb67d03580c80c773adaf3e38fd84d6b3211f6b69"}`))
	payload := &grequests.RequestOptions{
		JSON: j,
	}

	fmt.Println(payload)

	for {
		grequests.Post(url,payload)
	}
}