package main

import (
	"encoding/json"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

type Hello struct {
  Name string
  Msg string
}

func hello(c web.C, w http.ResponseWriter, r *http.Request)  {
  name := c.URLParams["name"]
  if name == ""{
    name = "gopher"
  }

  hello := &Hello{
    Name: name,
    Msg: "Hello",
  }
  encoder := json.NewEncoder(w)
  encoder.Encode(hello)
}

func main() {
  goji.Get("/hello/:name", hello)
  goji.Serve()
}
