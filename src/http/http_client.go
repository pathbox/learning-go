package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func httpGet() {
	resp, err := http.Get("http://www.baidu.com")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}

func httpPost() {
	resp, err := http.Post("http://www.baidu.com",
		"application/x-www-form-urlencodeed",
		strings.NewReader("name=cjb"))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}

func httpPostForm() {
	resp, err := http.PostForm("http://www.baidu.com",
		url.Values{"key": {"Value"}, "id": {"123"}})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))

}

func httpDo() {
	client = &http.Client{}

	req, err := http.NewRequest("POST", "http://www.baidu.com", strings.NewReader("name=cjb"))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "name=anny")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}

func main() {
	httpGet()
	httpPost()
	httpDo()
}
