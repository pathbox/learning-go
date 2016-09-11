package main

import (
	"fmt"

	"github.com/Jeffail/gabs"
)

func main() {
  jsonObj := gabs.New()
  jsonObj.Array("foo", "array")
  // or .ArrayP("foo.array")

  jsonObj.ArrayAppend(10, "foo", "array")
  jsonObj.ArrayAppend(20, "foo", "array")
  jsonObj.ArrayAppend(30, "foo", "array")

  fmt.Println(jsonObj.StringIndent("", "  "))
}
