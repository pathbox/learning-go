package main

import (
	"bytes"
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

	url := "http://127.0.0.1:9090" // server is ./post/server.go

	rp := &ReqParams{
		Action:            "PostJsonData",
		Backend:           "Backend",
		ResourceType:      1,
		TopOrganizationId: 2,
		ResourceId:        []string{"aaa", "bbb", "ccc"},
	}
	rpJSON, _ := json.Marshal(rp)
	payload := bytes.NewBuffer(rpJSON)
	request, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Errorf("Request internal_api error: %s", err)
	}

	request.Header.Set("Content-Type", "application/json") // It is must be set

	client := NewHTTPClient()
	resp, err := client.Do(request)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	res, _ := ioutil.ReadAll(resp.Body)
	// io.Copy(ioutil.Discard, resp.Body)
	log.Printf("response: %s", res)
	log.Print("It is Done")
}

// 自己定义一个params struct，然后用json序列化称为payload进行传递
type ReqParams struct {
	Action            string   `json:"Action"`
	Backend           string   `json:"Backend"`
	ResourceType      int      `json:"ResourceType"`
	TopOrganizationId int      `json:"TopOrganizationId"`
	ResourceId        []string `json:"ResourceId"`
}

func main() {
	DoPost()
}
