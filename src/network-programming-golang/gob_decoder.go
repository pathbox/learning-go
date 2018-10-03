package main

import (
	"fmt"

	"os"

	"encoding/gob"
)

type Person struct {
	Name

	Name

	Email []Email
}

type Name struct {
	Family

	string

	Personal string
}

type Email struct {
	Kind

	string

	Address string
}

func (p Person) String() string {

	s := p.Name.Personal + " " + p.Name.Family

	for _, v := range p.Email {

		s += "\n" + v.Kind + ": " + v.Address

	}

	return s

}

func main() {

	var person Person

	loadGob("person.gob", &person)

	fmt.Println("Person", person.String())

}

func loadGob(fileName string, key interface{}) {

	inFile, err := os.Open(fileName)

	checkError(err)

	decoder := gob.NewDecoder(inFile)

	err = decoder.Decode(key)

	checkError(err)

	inFile.Close()

}

func checkError(err error) {

	if err != nil {

		fmt.Println("Fatal error ", err.Error())

		os.Exit(1)

	}

}
