package main

import (
	"fmt"

	"github.com/thedevsaddam/gojsonq"
)

const json = `{"name":{"first":"Tom","last":"Hanks"},"age":61}`

func main() {

	j := gojsonq.New().JSONString(json)
	fmt.Println("json:", j)
	r1 := j.Pluck("name")
	fmt.Println("r1:", r1)
	r2 := j.Find("name")
	fmt.Println("r2:", r2)

	// r3 := j.Find("name.first") // it get nil, doesn't work
	r3 := gojsonq.New().JSONString(json).Find("name.first")
	fmt.Println("r3:", r3)
}
