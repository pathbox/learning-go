package main

import (
	"fmt"

	"github.com/Jeffail/gabs"
)

func main() {
	jsonParsed, _ := gabs.ParseJSON([]byte(`{"object":{ "first": 1, "second": 2, "third": 3 }}`))

// S is shorthand for search
  children, _ := jsonParsed.S("object").ChildrenMap()
  for key, child := range children {
    fmt.Printf("key: %v, value: %v\n", key, child.Data().(float64))
  }
}
