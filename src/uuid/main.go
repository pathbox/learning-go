package main

import (
	"fmt"

	"github.com/satori/go.uuid"
)

func main() {
	// Creating UUID Version 1
	u1 := uuid.NewV1()
	fmt.Printf("UUIDv1: %s\n", u1)

	// Creating UUID Version 2
	u2 := uuid.NewV2(uuid.DomainPerson)
	fmt.Printf("UUIDv2: %s\n", u2)

	// Creating UUID Version 3
	u3 := uuid.NewV3(uuid.NamespaceDNS, "www.example.com")
	fmt.Printf("UUIDv3: %s\n", u3)

	// Creating UUID Version 4
	u4 := uuid.NewV4()
	fmt.Printf("UUIDv4: %s\n", u4)

	// Creating UUID Version 5
	u5 := uuid.NewV5(uuid.NamespaceDNS, "www.example.com")
	fmt.Printf("UUIDv5: %s\n", u5)

	// Parsing UUID from string input
	up, err := uuid.FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	if err != nil {
		fmt.Printf("Something gone wrong: %s\n", err)
	}
	fmt.Printf("Successfully parsed: %s\n", up)
}
