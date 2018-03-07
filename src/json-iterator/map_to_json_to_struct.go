package main

import (
	"fmt"

	"github.com/json-iterator/go"
)

type User struct {
	Name string
	Age  string
}

func main() {
	m := make(map[string]string)

	m["name"] = "Joe"
	m["age"] = "27"

	b, _ := jsoniter.Marshal(m)
	user := &User{}
	jsoniter.Unmarshal(b, user)

	fmt.Println(user.Name)
	fmt.Println(user.Age)
	fmt.Println(user)
}
