package ydict

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type Result struct {
	ErrorCode   int
	Query       string
	Translation []string
	Basic       *struct {
		Phonetic string
		Explians []string
	}

	Web []struct {
		Key   string
		Value []string
	}
}

type Client struct {
	BaseURL string
	Keyfrom string
	Key     string
}

var OnlineBaseURL string = "http://fanyi.youdao.com/"

func NewClient(baseURL, keyfrom, key string) *Client {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL + "/"
	}

	return &client{
		BaseURL: baseURL,
		Keyfrom: keyfrom,
		Key:     key,
	}
}

func NewOnlineClient(keyfrom, key string) *Client {
	return NewClient(OnlineBaseURL, keyfrom, key)
}

type result struct {
	ErrorCode   int      `json:"errorCode"`
	Query       string   `json:"query"`
	Translation []string `json:"translation"`
	Basic       *struct {
		Phonetic string   `json:"phonetic"`
		Explains []string `json:"explains"`
	} `json:"basic"`
	Web []struct {
		Key   string   `json:"key"`
		Value []string `json:"value"`
	} `json:"web"`
}

func (r *result) asResult() *Result {
	res := &Result{
		ErrorCode:   r.ErrorCode,
		Query:       r.Query,
		Translation: r.Translation,
	}

	if r.Basic != nil {
		res.Basic = &struct {
			Phonetic string
			Explains []string
		}{
			Phonetic: r.Basic.Phonetic,
			Explains: r.Basic.Explains,
		}
	}

	// copy Web field
	res.Web = make([]struct {
		Key   string
		Value []string
	}, len(r.Web))
	for i := range r.Web {
		res.Web[i].Key = r.Web[i].Key
		res.Web[i].Value = r.Web[i].Value
	}

	return res
}

func (c *Client) Query(q string) (*Result, error) {
	return c.QueryHttp(http.DefaultClient, q)
}

func (c *Client) QueryHttp(httpClient *http.Client, q string) (*Result, error) {
	requestURL := fmt.Sprintf(
		"%sopenapi.do?keyfrom=%s&key=%s&type=data&doctype=json&version=1.1&q=%s",
		c.BaseURL, template.URLQueryEscaper(c.Keyfrom),
		template.URLQueryEscaper(c.Key), template.URLQueryEscaper(q))
	fmt.Println(requestURL)

	resp, err := httpClient.Get(requestURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var res result
	err = dec.Decode(&res)
	if err != nil {
		return nil, err
	}

	return res.asResult(), nil
}
