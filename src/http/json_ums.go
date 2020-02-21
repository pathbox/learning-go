package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	json "github.com/json-iterator/go"

	"github.com/labstack/gommon/log"
)

func NewHTTPClient() *http.Client {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
	}

	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}
	return client
}

func DoPost() {

	url := "http://xxx:4012" // server is ./post/server.go

	rp := map[string]interface{}{
		"Content":      "您的验证码是: 817723",
		"SendType":     "SMS",
		"Target":       "phonenumber",
		"request_uuid": "333",
		"SessionNo":    "6376251a-6f7a-4a3b-a9ce-08d7a963bf3e",
		// "ChannelId":    1,
	}
	rpJSON, _ := json.Marshal(rp)
	payload := bytes.NewBuffer(rpJSON)
	request, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Errorf("Request internal_api error: %s", err)
	}

	// request.Header.Set("Content-Type", "application/json") // It is must be set

	client := NewHTTPClient()
	resp, err := client.Do(request)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	res, _ := ioutil.ReadAll(resp.Body)
	// io.Copy(ioutil.Discard, resp.Body)
	// log.Printf("response: %s", res)
	fmt.Println("res", string(res))
	// log.Print("It is Done")
}

// 自己定义一个params struct，然后用json序列化称为payload进行传递
type ReqParams struct {
	Action       string `json:"Action"`
	AuthKey      string `json:"AuthKey"`
	Content      string `json:"Content"`
	SendType     string `json:"SendType"`
	Target       string `json:"Target"`
	PublicKey    string `json:"PublicKey"`
	Signature    string `json:"Signature"`
	request_uuid string `json:"request_uuid"`
	SessionNo    string `json:"SessionNo"`
}

func main() {
	DoPost()
}
