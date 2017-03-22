package main

import (
	"fmt"
)

type People struct {
	group   string
	*Person //  get all Person struct methods
	city    string
}

type Person struct{}

func main() {
	people := &People{}
	people.sayHello()
}

func (p *Person) sayHello() {
	fmt.Println("Hello World")
}
