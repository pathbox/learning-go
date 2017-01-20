package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

func main() {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSmFtZXMiLCJuYmYiOjE0NDQ0Nzg0MDB9.eCfz950i2QqveRAz72oK_nszEp1R31hoMtGbzm63dd0"

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		hmacSampleSecret := []byte("my secret")
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["name"], claims["nbf"])
	} else {
		fmt.Println(err)
	}
}
