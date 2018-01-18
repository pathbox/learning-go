package main

import (
	"github.com/valyala/fasthttp"
)

func DoRequest(url string) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.SetBodyString("p=q")
	req.Header.Add("User-Agent", "Test-Agent")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")
	println(req.Header.String())

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		println("Error:", err.Error())
	} else {
		bodyBytes := resp.Body()
		println(string(bodyBytes))
	}

	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
}

func main() {
	url := ""
	DoRequest(url)
}
