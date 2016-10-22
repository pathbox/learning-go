package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

type transport struct {
  http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
  resp, err := t.RoundTripper.RoundTrip(req)
  if err != nil {
    return nil, err
  }
  err = resp.Body.Close()
  if err != nil {
    return nil, err
  }

  b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)
  body := ioutil.NopCloser(bytes.NewReader(b))
  resp.Body = body
  resp.ContentLength = int64(len(b))
  resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
  resp.Header.Set("foo", "bar")
  return resp, nil
}

func main() {
  target, err := url.Parse("http://www.hao123.com")
  if err != nil {
    panic(err)
  }
  proxy := httputil.NewSingleHostReverseProxy(target)
  proxy.Transport = &transport{http.DefaultTransport}

  http.Handle("/", proxy)
  log.Fatal(http.ListenAndServe(":9090", nil))
}
