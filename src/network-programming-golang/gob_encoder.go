package main

import (
	"encoding/gob"
	"fmt"
	"os"
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

func main() {
	person := Person{

		Name: Name{Family: "Newmarch", Personal: "Jan"},

		Email: []Email{Email{Kind: "home", Address: "jan@newmarch.name"},

			Email{Kind: "work", Address: "j.newmarch@boxhill.edu.au"}}}

	saveGob("person.gob", person)
}
func checkError(err error) {

	if err != nil {

		fmt.Println("Fatal error ", err.Error())

		os.Exit(1)

	}

}
func saveGob(fileName string, key interface{}) {
	outFile, err := os.Create(fileName)
	checkError(err)
	encoder := gob.NewEncoder(outFile)
	err = encoder.Encode(key)
	checkError(err)
	defer outFile.Close()
}
