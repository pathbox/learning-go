package main

import (
	"fmt"
	"github.com/franela/goreq"
	"net/url" //don't forget to import net/url
)

func main() {
	values := url.Values{}
	values.Add("key", "value")
	values.Add("mostafa", "dahab")
	req, _ := goreq.Request{
		Uri:         "https://httpbin.org/post",
		Method:      "POST",
		Body:        values.Encode(),
		ContentType: "application/x-www-form-urlencoded; charset=UTF-8",
	}.Do()
	html, _ := req.Body.ToString()
	fmt.Println(html)
}
