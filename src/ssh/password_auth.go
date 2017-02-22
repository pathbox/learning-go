package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
)

func main() {
	fmt.Println("Here ssh auth start")

	config := &ssh.ClientConfig{ // one: config config
		User: "user",
		Auth: []ssh.AuthMethod{
			ssh.Password("123456654321"),
		},
	}

	client, err := ssh.Dial("tcp", "server_address_or_ip:22", config) // two: Dial conn to create client
	if err != nil {
		panic(err.Error())
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession() // three: NewSession()
	if err != nil {
		panic(err)
	}

	defer session.Close() // last: Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil { // four: sessoin.Run
		panic(err)
	}
	fmt.Println(b.String())
}
