package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func main() {
	baseURL := "http://127.0.0.1:8081/"

	usernamePwd := url.Values{}
	usernamePwd.Set("username", "suraj")
	usernamePwd.Set("password", "suraj")

	resp, err := http.PostForm(baseURL+"api/get-token/", usernamePwd)
	if err != nil {
		fmt.Println("Is the server running?")
		os.Exit(1)
	} else {
		fmt.Println("response received")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body")
	} else {
		fmt.Println("Token received")
	}
	token := string(body)

	client := &http.Client{}
	req, err := http.NewRequest("GET", baseURL+"api/get-task/", nil)

	if err != nil {
		fmt.Println("Unable to form a GET /api/get-task/")
	}

	req.Header.Add("Token", token)
	resp, err = client.Do(req)
	if (err != nil) || (resp.StatusCode != 200) {
		fmt.Println("Something went wrong in the getting a response")
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
