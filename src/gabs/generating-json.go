package main

import (
	"fmt"

	"github.com/Jeffail/gabs"
)

func main() {
  jsonObj := gabs.New()
  // or gabs.Consume(jsonObject) to work on an existing map[string]interface{}
  jsonObj.Set(10, "outter", "inner", "value")
  jsonObj.SetP(20, "outter.inner.value2")
  jsonObj.Set(30, "outter", "inner2", "value3")
  jsonObj.SetP("{'a':1, 'b':2}", "outter.inner.value4")

  fmt.Println(jsonObj.String())

  fmt.Println(jsonObj.StringIndent("", "  "))

}
