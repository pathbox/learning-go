res, err := goreq.Request{
    Method: "GET",
    Proxy: "http://myproxy:myproxyport",
    Uri: "http://www.google.com",
}.Do()
// Proxy basic auth is also supported

res, err := goreq.Request{
    Method: "GET",
    Proxy: "http://user:pass@myproxy:myproxyport",
    Uri: "http://www.google.com",
}.Do()