package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func main() {
	var transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		MaxIdleConns:          100, // KeepAlive 连接保持数量
		MaxIdleConnsPerHost:   100, // 对每个Host KeepAlive 连接保持数量 默认是2
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	var client = &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	res, _ := client.Get("http://www.baidu.com")

	defer res.Body.Close()

	if res.StatusCode == 200 {
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(body))
	}

	// n, err := io.Copy(ioutil.Discard, resp.Body)
}
