package main

import (
	"fmt"
	"github.com/franela/goreq"
)

type Item struct {
	TheLimit  int    `url:"the_limit"`
	TheSkip   string `url:"the_skip,omitempty"`
	TheFields string `url:"-"`
}

func main() {
	// item := Item{
	// 	Limit:  3,
	// 	Skip:   5,
	// 	Fields: "Value",
	// }
	res, err := goreq.Request{Uri: "http://www.baidu.com"}.Do()
	// item := url.Values{}
	// item.Set("Limit", 3)
	// item.Add("Field", "somefield")
	// item.Add("Field", "someotherfield")

	// res, err := goreq.Request{Uri: "https://www.baidu.com", QueryString: item}.Do()

	if err != nil {
		return
	}

	fmt.Println(res.Body)
}
