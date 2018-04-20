package main

import "fmt"
import "encoding/json"

type Detail struct {
	Time string `json:"time"`
	Ops  string `json:"ops"`
}

type Version struct {
	Last   string   `json:"last"`
	Max    int      `json:"max"`
	Detail []Detail `json:"detail"`
}

type Desc struct {
	Date    string  `json:"date"`
	Content string  `json:"content"`
	Brief   string  `json:"brief"`
	Keyword string  `json:"keyword"`
	Version Version `json:"version"`
}

type Response struct {
	Url     string   `json:"url"`
	Title   string   `json:"title"`
	Email   string   `json:"email"`
	Admin   string   `json:"admin"`
	Address []string `json:"address"`
	Article []Desc   `json:"article"`
}

func main() {
	body := `
  {
      "url": "http://xiaorui.cc",
      "title": "golang and python",
      "admin": "fengyun",
      "email": "rfyiamcool@163.com",
      "address": [
          "beijing",
          "qingdao"
      ],
      "article": [
          {
              "date": "2014",
              "content": "golang json push to redis server",
              "brief": "golang json",
              "keyword": "json",
              "version": {
                  "max": 3,
                  "last": "2016-03-11",
                  "detail": [
                      {
                          "time": "2016-03-12",
                          "ops": "add my email"
                      }
                     ]
                  }
          }
      ]
  }
  `

	var r Response
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		fmt.Printf("err was %v", err)
	}

	fmt.Println(r.Article[0].Version.Max)
	fmt.Println(r.Article[0].Version.Detail[0].Ops)
	fmt.Println(r)

}
