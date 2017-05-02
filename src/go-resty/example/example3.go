package main

import (
	"gopkg.in/resty.v0"
	"log"
)

type AuthSuccess struct{}

func main() {
	// POST JSON string
	// No need to set content type, if you have client level setting
	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"username":"testuser", "password":"testpass"}`).
		SetResult(&AuthSuccess{}).
		Post("https://myapp.com/login")

		// POST []byte array
	// No need to set content type, if you have client level setting
	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody([]byte(`{"username":"testuser", "password":"testpass"}`)).
		SetResult(&AuthSuccess{}). // or SetResult(AuthSuccess{}).
		Post("https://myapp.com/login")

		// POST Struct, default is JSON content type. No need to set one
	resp, err := resty.R().
		SetBody(User{Username: "testuser", Password: "testpass"}).
		SetResult(&AuthSuccess{}). // or SetResult(AuthSuccess{}).
		SetError(&AuthError{}).    // or SetError(AuthError{}).
		Post("https://myapp.com/login")

		// POST Map, default is JSON content type. No need to set one
	resp, err := resty.R().
		SetBody(map[string]interface{}{"username": "testuser", "password": "testpass"}).
		SetResult(&AuthSuccess{}). // or SetResult(AuthSuccess{}).
		SetError(&AuthError{}).    // or SetError(AuthError{}).
		Post("https://myapp.com/login")

		// POST of raw bytes for file upload. For example: upload file to Dropbox
	fileBytes, _ := ioutil.ReadFile("/Users/jeeva/mydocument.pdf")

	resp, err := resty.R().
		SetBody(fileBytes).
		SetContentLength(true).
		SetAuthToken("my_token").
		SetError(&DropboxError{}).
		Post("https://content.dropboxapi.com/1/files_put/auto/resty/mydocument.pdf")

	log.Printf("\nError: %v", err)
	log.Printf("\nResponse Status Code: %v", resp.StatusCode())
	log.Printf("\nResponse Status: %v", resp.Status())
	log.Printf("\nResponse Time: %v", resp.Time())
	log.Printf("\nResponse Recevied At: %v", resp.ReceivedAt())
	log.Printf("\nResponse Body: %v", resp) // or resp.String() or string(resp.Body())
}
