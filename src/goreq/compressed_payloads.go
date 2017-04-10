// #####Using gzip compression:

res, err := goreq.Request{
    Method: "POST",
    Uri: "http://www.google.com",
    Body: item,
    Compression: goreq.Gzip(),
}.Do()
// #####Using deflate/zlib compression:

res, err := goreq.Request{
    Method: "POST",
    Uri: "http://www.google.com",
    Body: item,
    Compression: goreq.Deflate(),
}.Do()

type Item struct {
    Id int
    Name string
}
res, err := goreq.Request{
    Method: "POST",
    Uri: "http://www.google.com",
    Body: item,
    Compression: goreq.Gzip(),
}.Do()
var item Item
res.Body.FromJsonTo(&item)