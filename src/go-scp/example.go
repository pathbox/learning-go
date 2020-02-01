package main

import (
	"fmt"
	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
	"os"
)

func main() {
	clientConfig, _ := auth.PrivateKey("username", "/path/to/rsa/key", ssh.InsecureIgnoreHostKey())

	client := scp.NewClient("example.com:22", &clientConfig)

	err := client.Connect()

	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)
		return
	}

	// Open a file
	f, _ := os.Open("/path/to/local/file")

	// Close client connection after the file has been copied
	defer client.Close()

	// Close the file after it has been copied
	defer f.Close()
	err = client.CopyFile(f, "/home/server/test.txt", "0655")
	if err != nil {
		fmt.Println("Error while copying file ", err)
	}
}