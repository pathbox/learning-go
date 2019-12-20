https://github.com/google/go-querystring

### go-querystring is Go library for encoding structs into URL query parameters.

```go
// Package query implements encoding of structs into URL query parameters.
//
// As a simple example:
//
// 	type Options struct {
// 		Query   string `url:"q"`
// 		ShowAll bool   `url:"all"`
// 		Page    int    `url:"page"`
// 	}
//
// 	opt := Options{ "foo", true, 2 }
// 	v, _ := query.Values(opt) // 1. Values get v
// 	fmt.Print(v.Encode()) // will output: "q=foo&all=true&page=2" 2. v.Encode() get the url query string
//
// The exact mapping between Go values and url.Values is described in the
```