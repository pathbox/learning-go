res, err := goreq.Request{
  Method:      "GET",
  Uri:         "http://www.google.com",
  Compression: goreq.Gzip(),
  ShowDebug:   true,
}.Do()
fmt.Println(res, err)