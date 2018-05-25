package postclientexample

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// "application/x-www-form-urlencoded"

func DoPost() {
	var r http.Request
	r.ParseForm()

	r.Form.Add("name", "Joe")
	r.Form.Add("age", "27")

	url := "127.0.0.1:9090/test/post"
	payload := strings.TrimSpace(r.Form.Encode())
	request, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// request.Header.Set("Content-Type", "application/json")
	// request.Header.Set("Connection", "Keep-Alive")

	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	io.Copy(ioutil.Discard, resp.Body)
}
