package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func randString() string {
	b := make([]byte, 10)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func main() {
	fmt.Println(randString())
}
