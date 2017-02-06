package main

import (
	"fmt"
	"github.com/elithrar/simple-scrypt"
	"log"
)

func main() {
	passwordFromForm := "prew8fid9hick6c"

	hash, err := scrypt.GenerateFromPassword([]byte(passwordFromForm), scrypt.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hash)

	err := scrypt.CompareHashAndPassword(hash, []byte(passwordFromForm))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(hash)
}
