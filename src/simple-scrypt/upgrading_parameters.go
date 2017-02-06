package main

import (
	"github.com/elithrar/simple-scrypt"
	"log"
)

func main() {
	current, err := scrypt.Cost(hash)
	if err != nil {
		log.Fatal(err)
	}

	slower := scrypt.Params{
		N:       32768,
		R:       8,
		P:       2,
		SaltLen: 16,
		DKLen:   32,
	}

	if !reflect.DeepEqual(current, slower) {
		// Re-generate the key with the slower parameters
		// here using scrypt.GenerateFromPassword
	}
}
