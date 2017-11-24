package main

import "fmt"
import "github.com/json-iterator/go"

func main() {
	b := []byte(`{"Name": "Wednesday", "Age": 6, "Parents": ["Gomez", "Morticia"]}`)
	fmt.Println("b", b)

	var f interface{}
	err := jsoniter.Unmarshal(b, &f)
	if err != nil {
		panic(err)
	}
	fmt.Println("f", f)
	m := f.(map[string]interface{})
	fmt.Println("Name", m["Name"])
}
