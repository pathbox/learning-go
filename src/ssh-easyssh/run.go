package main

import (
	"fmt"

	"github.com/hypersleep/easyssh"
)

func main() {
	// Create MakeConfig instance with remote username, server address and path to private key.

	key_file := "/home/user/.ssh/id_rsa"

	ssh := &easyssh.MakeConfig{
		User:   "webuser",
		Server: "udesk.test.dog",
		// Optional key or Password without either we try to contact your agent SOCKET
		//Password: "password",
		Key:  key_file,
		Port: "22",
	}

	stdout, stderr, done, err := ssh.Run("ps ax", 60)
	if err != nil {
		panic("Can't run remote command: " + err.Error())
	} else {
		fmt.Println("don is :", done, "stdout is :", stdout, ";   stderr is :", stderr)
	}

}
