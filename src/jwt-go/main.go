package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"log"
	"net/http"
	"strings"
	"time"
)

type MyCustomClaims struct {
	// This will hold a users username after authenticating.
	// Ignore `json:"username"` it's required by JSON
	Username string `json:"username"`

	// This will hold claims that are recommended having (Expiration, issuer)
	jwt.StandardClaims
}

func setToken(res http.ResponseWriter, req *http.Request) {
	// Expires the token and cookie in 24 hours
	expireToken := time.Now().Add(time.Hour * 24).Unix()
	expireCookie := time.Now().Add(time.Hour * 24)
	// We'll manually assign the claims but in production you'd insert values from a database
	claims := MyCustomClaims{
		"Curry",
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "example.com",
		},
	}
	// Create the token using your claims
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	signedToken, _ := token.SignedString([]byte("secret"))
	// This cookie will store the token on the client side
	cookie := http.Cookie{Name: "Auth", Value: signedToken, Expires: expireCookie, HttpOnly: true}
	http.SetCookie(res, &cookie)

	http.Redirect(res, req, "/profile", 301)
}

func validate(protectedPage http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		cookie, err := req.Cookie("Auth")
		if err != nil {
			http.NotFound(res, req)
			return
		}

		splitCookie := strings.Split(cookie.String(), "Auth=")
		log.Println("=================")
		log.Println(splitCookie)

		token, err := jwt.ParseWithClaims(splitCookie[1], &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})

		if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
			context.Set(req, "Claims", claims)
		} else {
			http.NotFound(res, req)
			return
		}
		protectedPage(res, req)
	})
}

func profile(res http.ResponseWriter, req *http.Request) {
	claims := context.Get(req, "Claims").(*MyCustomClaims)
	res.Write([]byte(claims.Username))
	context.Clear(req)
}

func homePage(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Home Page"))
}

func main() {
	http.HandleFunc("/profile", validate(profile))
	http.HandleFunc("/setToken", setToken)
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":9090", nil)
}
