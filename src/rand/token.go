package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func main() {
	token, err := GenerateRandomString(32)
	if err != nil {
		panic(err)
	}
	fmt.Println("Generate Token: ", token)
}
