package main

import (
	"fmt"

	"github.com/thedevsaddam/gojsonq"
)

func main() {
	jq := gojsonq.New().File("./sample-data.json")
	res := jq.From("vendor.items").Where("price", ">", 1200).OrWhere("id", "=", nil).Only("name", "price")

	fmt.Println(res)
}
