package main

import (
	"fmt"

	"github.com/Jeffail/gabs"
)

func main() {
	jsonParsed, _ := gabs.ParseJSON([]byte(`{"array":[ {"value1":1}, {"value":2}, {"value":3} ]}`))
	fmt.Println(jsonParsed.Path("array.value").String())
}
