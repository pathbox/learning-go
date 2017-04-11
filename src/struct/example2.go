package main

import (
	"./person"
	"fmt"
)

func main() {
	person := &person.Person{Name: "Curry", Age: 28}

	fmt.Println(person.Name)
}
