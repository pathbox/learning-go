package main

import (
	"fmt"
	"net/http"
)

func main() {
	client := &http.Client{}
	UpercaseHeaderRequest(client)
	LowercaseHeaderRequest(client)
}

// Headers中传 X-API-KEY
/*
   However, request.Header.Set(...) ends up calling CanonicalMIMEHeaderKey on the header key, in this case "x-api-key". This converts "x-api-key" to "X-Api-Key". Note that request.Header.Add("x-api-key", "somelongapikey2349208759283") does the same thing. While this is good behavior on Go’s part, we need this header to be all lowercase.
*/

func UpercaseHeaderRequest(client *http.Client) {
	// req, _ := http.NewRequest("GET", "https://someapi/somesource", nil)
	req, _ := http.NewRequest("GET", "https://www.baidu.com", nil)
	req.Header.Set("x-api-key", "somelongapikey1234567890")
	fmt.Println(req.Header)
	res, _ := client.Do(req)
	fmt.Println(res)

}

// Headers 中传 x-api-key
/*
So how can we set a header to be all lowercase? It turns out that request.Header is an alias of the type map[string][]string. Thus, we can set the header key as lowercase (or whatever we want) with the following code:
*/

func LowercaseHeaderRequest(client *http.Client) {
	req, _ := http.NewRequest("GET", "https://www.baidu.com", nil)
	req.Header["x-api-key"] = []string{"somelongapikey2349208759283"}
	fmt.Println(req.Header)
	res, _ := client.Do(req)
	fmt.Println(res)
}
