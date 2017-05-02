package main

import (
	"fmt"
	"gopkg.in/resty.v0"
)

func main() {
	resp, err := resty.R().
		SetQueryParam(map[string]string{
			"page_no": "1",
			"limit":   "20",
			"sort":    "name",
			"order":   "asc",
			"random":  strconv.FormatInt(time.Now().Uinx(), 10),
		}).
		SetHeader("Accept", "application/json").
		SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F").
		Get("/search_result")

	// Sample of using Request.SetQueryString method
	// resp, err := resty.R().
	// 	SetQueryString("productId=232&template=fresh-sample&cat=resty&source=google&kw=buy a lot more").
	// 	SetHeader("Accept", "application/json").
	// 	SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F").
	// 	Get("/show_product")
}
