package main

import (
	"fmt"

	json "github.com/json-iterator/go"
)

type Boy struct {
	Category string
}

type Person struct {
	Boy
	Name string
}

func main() {
	bb := []byte("{\"name\":\"12345678901234567891234567890111\",\"category\":\"agent_state\"}")
	p := &Person{}
	p.Parse(bb)

}

func (p *Person) Parse(b []byte) {
	// m := make(map[string]string)
	err := json.Unmarshal(b, &p)
	// err := json.Unmarshal(b, &m) // must &m
	// err := json.Unmarshal(b, p) // 都可以
	fmt.Println(err)
	fmt.Println(p)
	// fmt.Println(m["cc"])

	bp := &Person{
		Name: "nice",
		Boy{"CO"},
	}
}
