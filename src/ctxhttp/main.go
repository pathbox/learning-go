package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/Jeffail/gabs"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

func main() {
  ctx, cancel := content.WithTimeout(context.Background(), 3*time.Second)

  go func(){
    cancel()
  }()

  resp, err := ctxhttp.Get(ctx, nil, "http://httppin.org/get")
  if err != nil {
    panic(err)
  }

  resBody, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    panic(err)
  }

  jsonParsed, err := gabs.ParseJSON(resBody)
  if err != nil {
    panic(err)
  }

  fmt.Println(jsonParsed.StringIndent("", " "))
}
