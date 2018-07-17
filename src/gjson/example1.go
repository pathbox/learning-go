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

	r := gjson.Get(json, "friends")

	for _, i := range r.Array() {
		m := i.Map()
		fmt.Println(m["first"].String())
	}
}
