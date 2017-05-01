package main

import (
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/kr/pretty"
	"gopkg.in/jmcvetta/napping.v3"
	"log"
	"net/url"
	"time"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func main() {

	var username string
	fmt.Printf("Github username: ")
	_, err := fmt.Scanf("%s", &username)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Hello Here")
	fmt.Printf("github.com/howeyc/gopass")
	passwd, err := gopass.GetPasswd()
	if err != nil {
		log.Fatal(err)
	}

	//
	// Compose request
	//
	// http://developer.github.com/v3/oauth/#create-a-new-authorization
	//

	payload := struct { // payload  request params
		Scopes []string `json:"scopes"`
		Note   string   `json:"note"`
	}{
		Scopes: []string{"public_repo"},
		Note:   "testing Go napping" + time.Now().String(),
	}

	res := struct { // response  response struct
		Id        int
		Url       string
		Scopes    []string
		Token     string
		App       map[string]string
		Note      string
		NoteUrl   string `json:"note_url"`
		UpdatedAt string `json:"updated_at"`
		CreatedAt string `json:"created_at"`
	}{}

	//
	// Struct to hold error response
	//
	e := struct {
		Message string
		Errors  []struct {
			Resource string
			Field    string
			Code     string
		}
	}{}

	//
	// Setup HTTP Basic auth for this session (ONLY use this with SSL).  Auth
	// can also be configured on a per-request basis when using Send().
	//
	s := napping.Session{
		UserInfo: url.UserPassword(username, string(passwd)),
	}

	url := "https://api.github.com/authorizations"

	resp, err := s.Post(url, &payload, &res, &e)

	if err != nil {
		log.Fatal(err)
	}

	println("")
	if resp.Status() == 201 {
		fmt.Printf("Github auth token: %s\n\n", res.Token)
	} else {
		fmt.Println("Bad response status from Github server")
		fmt.Printf("\t Status:  %v\n", resp.Status())
		fmt.Printf("\t Message: %v\n", e.Message)
		fmt.Printf("\t Errors: %v\n", e.Message)
		pretty.Println(e.Errors)
	}

}
