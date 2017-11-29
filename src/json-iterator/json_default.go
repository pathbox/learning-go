package main

import "fmt"
import "github.com/json-iterator/go"

type User struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
	Code string `json:"code,omitempty"`
}

func main() {
	b := []byte(`{"Name": "Wednesday", "Age": 6}`)
	fmt.Println("b", b)

	user := &User{}
	err := jsoniter.Unmarshal(b, user)
	if err != nil {
		panic(err)
	}

	fmt.Println("user", user)

	if user.Code == "" {

		fmt.Println("Code", user.Code)
	} else {
		fmt.Println("none")
	}

}
