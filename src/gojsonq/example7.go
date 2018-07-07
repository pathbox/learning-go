package main

import (
	"fmt"

	"github.com/thedevsaddam/gojsonq"
)

func main() {
	jq := gojsonq.New().File("./sample-data.json")
	res := jq.From("vendor.items").SortBy("price").Get()

	fmt.Println(res)
}
