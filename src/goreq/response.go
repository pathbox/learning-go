res, err := goreq.Request {
  Uri: "http://www.google.com",
  Timeout: 500 * time.Millisecond,
}

res.Uri
res.StatusCode
res.Body
res.Body.ToString()
res.Header.Get("Content-Type")