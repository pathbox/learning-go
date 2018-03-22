package postclientexample

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// body全部二进制数据流进行post
func DBytesPost(url string, data []byte) ([]byte, error) {
	payload := bytes.NewReader(data)

	request, err := http.NewRequest("POST", url, payload)
	if err != nil {
		panic(err)
	}

	request.Header.Set("Connection", "Keep-Alive")
	var resp *http.Response

	resp, err = http.DefaultClient.Do(request)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	// io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		panic(err)
	}

	return b, err
}
