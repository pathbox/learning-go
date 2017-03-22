package main

import (
	"fmt"
)

type Persion struct {
	name    string
	age     string
	country string
	flag    int
}

func main() {
	person1 := &Persion{name: "Carry", age: "20", flag: 1}
	person2 := &Persion{name: "John", age: "20", flag: 2}
	person3 := &Persion{name: "Nansi", country: "USA", flag: 3}
	fmt.Println(person1.name)
	f := showName(person1.name)
	person1.introduce(f, "Have a nice day")
	f = showCountry(person2.age)
	person2.introduce(f, "Hey man pass the ball")
	f = showCountry(person3.country)
	person3.introduce(f, "Call me")

}

func (p *Persion) introduce(f func(v string), va string) {
	f(va)
}

func showName(v string) string {
	fmt.Println("My name is: ", v)
}

func showAge(v string) string {
	fmt.Println("My age is: ", v)
}

func showCountry(v string) string {
	fmt.Println("My country is: ", v)
}
