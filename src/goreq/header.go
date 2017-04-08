req := goreq.Request{ Uri: "http://www.google.com" }

req.AddHeader("X-Custom", "somevalue")

req.Do()

res, err = goreq.Request{ Uri: "http://www.google.com" }.WithHeader("X-Custom", "somevalue").Do()