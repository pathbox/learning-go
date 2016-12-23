package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func main() {
	u, _ := url.Parse("http://www.baidu.com")
	q := u.Query()
	q.Set("username", "user")
	q.Set("password", "passwd")
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
		return
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("%s", result)
}
