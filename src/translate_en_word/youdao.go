package main

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

var youdaoURL = "http://dict.youdao.com/dictvoice?type=%d&audio=%s"
var youdaomURL = "http://m.youdao.com/dict?le=eng&q="

type youdao struct{}

func (e youdao) audio(word string, us bool) (mp3, ipa, def string, err error) {
	var u string
	if us {
		u = fmt.Sprintf(youdaoURL, 2, word)
	} else {
		u = fmt.Sprintf(youdaoURL, 1, word)
	}

	mu := youdaomURL + word
	req, _ := http.NewRequest(http.MethodGet, mu, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 10_1_1 like Mac OS X) AppleWebKit/602.2.14 (KHTML, like Gecko)")
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("failed to get ipa from youdao: %v\n", err)
		return mp3, ipa, def, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("failed to read response from cambridge: %v\n", err)
		return u, "", "", err
	}

	if us {
		ipa = doc.Find("span.phonetic").First().Text()

	} else { //uk
		ipa = doc.Find("span.phonetic").First().Text()
	}

	doc.Find("#bd ul").Children().Each(func(_ int, s *goquery.Selection) {
		def = def + s.Text() + "\n"
	})

	return u, ipa, def, nil
}
