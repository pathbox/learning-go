package main

import (
	"fmt"

	"github.com/Jeffail/gabs"
)

func main() {
	jsonObj := gabs.New()

	// Create an array with the length of 3
	jsonObj.ArrayOfSize(3, "foo")

	jsonObj.S("foo").SetIndex("test1", 0)
	jsonObj.S("foo").SetIndex("test2", 1)

	// Create an embedded array with the length of 3
	jsonObj.S("foo").ArrayOfSizeI(3, 2)

	jsonObj.S("foo").Index(2).SetIndex(1, 0)
	jsonObj.S("foo").Index(2).SetIndex(2, 1)
	jsonObj.S("foo").Index(2).SetIndex(3, 2)

	fmt.Println(jsonObj.StringIndent("", "  "))
}
