package main

import (
	"fmt"

	"github.com/Jeffail/gabs"
)

func main() {
  jsonParseObj, _ := gabs.ParseJSON([]byte(`{
    "outter":{
        "inner":{
            "value1":10,
            "value2":22
        },
        "alsoInner":{
            "value1":20
        }
    }
  }`))
  jsonOutput := jsonParseObj.Search("outter").String()
  fmt.Println(jsonOutput)
}
