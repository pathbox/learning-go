package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gorilla/securecookie"
	"github.com/joho/godotenv"
	"github.com/knq/jwt"
	_ "github.com/lib/pq"
	"github.com/matryer/way"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// https://github.com/nicolasparada/go-github-oauth-demo

var origin *url.URL
var db *sql.DB
var githubOAuthConfig oauth2.Config
var sc *securecookie.SecureCookie
var signer jwt.Signer

func main() {
	godotenv.Load()

	port := intEnv("PORT", 3000)
	originString := env("ORIGIN", fmt.Sprintf("http://localhost:%d/", port))
	databaseURL := env("DATABASE_URL", "postgresql://root@127.0.0.1:26257/github_oauth?sslmode=disable")
	githubClientID := os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	hashKey := env("HASH_KEY", "secret")
	jwtSecret := env("JWT_SECRET", "secret")

	var err error
	if origin, err = url.Parse(originString); err != nil || !origin.IsAbs() {
		log.Fatal("invalid origin")
		return
	}

	if githubClientID == "" || githubClientSecret == "" {
		log.Fatalf("remember to set both $GITHUB_CLIENT_ID and $GITHUB_CLIENT_SECRET")
		return
	}

	if i, err := strconv.Atoi(origin.Port()); err == nil {
		port = i
	}

	if db, err = sql.Open("postgres", databaseURL); err != nil {
		log.Fatalf("could not open database connection: %v\n", err)
		return
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("could not ping to db: %v\n", err)
		return
	}

	githubRedirectURL := *origin
	githubRedirectURL.Path = "/api/oauth/github/callback"
	githubOAuthCOnfig = oauth2.Config{
		ClientID:     githubClientID,
		ClientSecret: githubClientSecret,
		Endpoint:     github.Endpoint,
		RedirectURL:  githubRedirectURL.String(),
		Scopes:       []string{"read:user"},
	}

	sc = securecookie.New([]byte(hashKey), nil).MaxAge(0)

	signer, err = jwt.HS256.New([]byte(jwtSecret))

	if err != nil {
		log.Fatalf("could not create JWT signer: %v\n", err)
		return
	}

	router := way.NewRouter()
	router.HandleFunc("GET", "/api/oauth/github", githubOAuthStart)

	router.HandleFunc("GET", "/api/oauth/github/callback", githubOAuthCallback)
	router.HandleFunc("GET", "/api/auth_user", guard(getAuthUser))
	router.HandleFunc("*", "/api/...", http.NotFound)
	router.Handle("GET", "/...", http.FileServer(http.Dir("static")))

	log.Printf("accepting connections on port %d\n", port)
	log.Printf("starting server at %s\n", origin.String())
	addr := fmt.Sprintf(":%d", port)

	if err = http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}

}

func env(key, fallbackValue string) string {
	v, ok := os.LookupEnv(key)

	if !ok {
		return fallbackValue
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return fallbackValue
	}
	return i
}
