package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	usernamePwd := url.Values{}
	usernamePwd.Set("username", "admin")
	usernamePwd.Set("password", "admin")
	fmt.Println(usernamePwd)

	resp, err := http.PostForm("http://127.0.0.1:8081/api/get-token/", usernamePwd)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}
