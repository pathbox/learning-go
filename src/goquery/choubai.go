package main

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	doc, err := goquery.NewDocument("http://www.qiushibaike.com")
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".content").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
	})
}
