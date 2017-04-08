package main

import (
	"fmt"
	"github.com/franela/goreq"
)

func main() {

	type Item struct {
		Id   int
		Name string
	}

	item := Item{Id: 1111, Name: "foobar"}

	res, err := goreq.Request{
		Method: "POST",
		Uri:    "http://www.google.com",
		Body:   item,
	}.Do()

	res1, err := goreq.Request{
		Uri:         "http://www.google.com",
		Host:        "foobar.com",
		Accept:      "application/json",
		ContentType: "application/json",
		UserAgent:   "goreq",
	}.Do()

}
