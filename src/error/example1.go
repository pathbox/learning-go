package main

import (
	"fmt"
	"net/http"
)

//向外层传递异常
func getLink() (*http.Response, error) {
	content, err := http.Get("www.baidu.com")
	if err != nil {
		return nil, err
	}
	return content, nil
}

func main() {
	resp, err := getLink()
	if err != nil {
		return
	}
	fmt.Println(resp.Body)
}
