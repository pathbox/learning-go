package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	//"github.com/360EntSecGroup-Skylar/excelize"
)

func ParsePage(url string) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review item
	doc.Find("body .main-content .basic-info.cmn-clearfix dl.basicInfo-block").Each(func(i int, s *goquery.Selection) {

		// aString, _ := s.Find("td").Find("a").Attr("href")
		// title, _ := s.Find("td").Find("a").Attr("title")
		// date := s.Find(".time").Text()

		// name := s.Find("dd.basicInfo-item.value").Text()
		name := s.Find("dt.basicInfo-item.name").Text()
		name = strings.TrimSpace(name)
		//if strings.Contains(name,"开发商") {
		//	fmt.Println("xx",name)
		//	v := s.Find("dd.basicInfo-item.value").Text()
		//	fmt.Println(v)
		//}
		fmt.Println("=",i,name)
	})
}

func main() {
	appName := "弹弹堂S"
	baikeUrl := fmt.Sprintf("https://baike.baidu.com/item/%s", appName)
	ParsePage(baikeUrl)
}
