package main

import (
	"flag"
	"fmt"
	"net/url"
	"strings"
)

type UrlFlag struct {
	urls []*url.URL
}

func (arr *UrlFlag) Urls() []*url.URL {
	return arr.urls
}

func (arr *UrlFlag) String() string {
	return fmt.Sprint(arr.urls)
}

func (arr *UrlFlag) Set(value string) error {
	if len(arr.urls) > 0 {
		return fmt.Errorf("The url flag is already set")
	}

	urls := strings.Split(value, ",")
	for _, item := range urls {
		if parsedUrl, err := url.Parse(item); err != nil {
			return err
		} else {
			arr.urls = append(arr.urls, parsedUrl)
		}
	}
	return nil
}

func main() {
	var arg UrlFlag
	flag.Var(&arg, "url", "URL comma-separated list")
	flag.Parse()

	for _, item := range arg.Urls() {
		fmt.Printf("scheme: %s url: %s path: %s\n", item.Scheme, item.Host, item.Path)
	}
}
