package main

import (
	"fmt"

	"github.com/tidwall/gjson"
)

func main() {
	json := `{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44},
    {"first": "Roger", "last": "Craig", "age": 68},
    {"first": "Jane", "last": "Murphy", "age": 47}
  ]
	}`

	t := gjson.Parse(json)
	fmt.Println(t.Get("friends").IsArray())
	fmt.Println(t.Get("name").IsObject())
	fmt.Println(t.Get("friends").IsObject())
	fmt.Println(t.Get("age").IsArray())
	fmt.Println(t.Get("age").IsObject())

}
