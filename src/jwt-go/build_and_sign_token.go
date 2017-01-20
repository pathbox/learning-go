package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func main() {
	claims := jwt.MapClaims{
		"name": "James",
		"nbf":  time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	hmacSampleSecret := []byte("my secret")
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tokenString)
}

// Output: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSmFtZXMiLCJuYmYiOjE0NDQ0Nzg0MDB9.eCfz950i2QqveRAz72oK_nszEp1R31hoMtGbzm63dd0
