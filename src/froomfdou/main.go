package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"./tool"

	"github.com/PuerkitoBio/goquery"
)

func ParsePage(url, keyword string) {
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

	// Find the review items
	doc.Find("#wrapper #content .article table.olt tbody tr").Each(func(i int, s *goquery.Selection) {

		aString, _ := s.Find("td").Find("a").Attr("href")
		title, _ := s.Find("td").Find("a").Attr("title")
		date := s.Find(".time").Text()

		if strings.Contains(title, keyword) {
			fmt.Printf("View %d:%s\n", i, title)
			fmt.Println(aString)
			fmt.Println(date)
		}
	})
}

func main() {

	// doubanUrl := "https://www.douban.com/group/shanghaizufang/discussion?start="
	// var url string
	// keyword := "宁国路"
	// for i := 0; i < 50; i++ {
	// 	time.Sleep(100 * time.Millisecond)
	// 	offset := strconv.Itoa(i * 25)
	// 	url = doubanUrl + offset
	// 	ParsePage(url, keyword)
	// }

	info := tool.NewEmail("luodaokai@udesk.cn", "test", "This is a test")
	tool.SendEmail(info)
}
