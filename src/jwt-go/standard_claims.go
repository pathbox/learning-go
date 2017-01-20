package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

func main() {
	mySigningKey := []byte("my secret")
	claims := &jwt.StandardClaims{
		ExpiresAt: 15000,
		Issuer:    "James",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	fmt.Println(err)
	fmt.Println(ss)
}
