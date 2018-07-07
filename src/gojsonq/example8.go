package main

import (
	"fmt"

	"github.com/thedevsaddam/gojsonq"
)

func main() {
	jq := gojsonq.New().File("./users.json")
	res := jq.From("users").WhereEqual("name.first", "John").Get()

	fmt.Println(res)
}
