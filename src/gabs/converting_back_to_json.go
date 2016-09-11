package main

import (
	"fmt"

	"github.com/Jeffail/gabs"
)

func main() {
  jsonParseObj, _ := gabs.ParseJSON([]byte(`{
    "outter":{
      "inner":{
        "value1": 10,
        "value2": 22
    },
    "alsoInner":{
      "value": 30
    }
  }
 }`))
 jsonOutput := jsonParseObj.StringIndent("", " ")
 fmt.Println(jsonOutput)

}
//
// {
//  "outter": {
//   "alsoInner": {
//    "value": 30
//   },
//   "inner": {
//    "value1": 10,
//    "value2": 22
//   }
//  }
// }
