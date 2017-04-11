res, err := goreq.Request {
  Uri: "http://www.google.com",
  Timeout: 500 * time.Millisecond,
}