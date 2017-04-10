res, err := goreq.Request{
  Uri: "http://.google.com"
}.
WithCookie(&http.Cookie{Name: "c1", Value: "v1"}).
Do()
