package main

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	url := "http://api.stackoverflow.com/1.1/tags?pagesize=100&page=1"
	res, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", string(body))

	js, err := simplejson.NewJson(body)
	if err != nil {
		log.Fatal(err)
	}

	total := js.Get("total").MustInt()

	fmt.Printf("Total: %s", total)
}
