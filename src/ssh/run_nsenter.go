package main

import (
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func main() {
	privateKey := privateKey()

	// Create client config
	config := &ssh.ClientConfig{
		User: "core",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(privateKey),
		},
	}

	conn, err := ssh.Dial("tcp", "myhost.com:22", config)
	if err != nil {
		log.Fatalf("unable to connect: %s", err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("unable to create session: %s", err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	// Stuff for interactive shell
	// Set up terminal modes
	//modes := ssh.TerminalModes{
	//  ssh.ECHO:          1,     // enable echoing
	//  ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
	//  ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	//}
	// Request pseudo terminal
	//if err := session.RequestPty("xterm-256color", 80, 40, modes); err != nil {
	//  log.Fatalf("request for pseudo terminal failed: %s", err)
	//}
	// Start remote shell
	//if err := session.Shell(); err != nil {
	//  log.Fatalf("failed to start shell: %s", err)
	//}

	if err := session.Run("sudo nsenter --target 2202 --mount --uts --ipc --net --pid"); err != nil {
		panic("Failed to run: " + err.Error())
	}
	session.Wait()
}

func privateKey() ssh.Signer {
	buf, err := ioutil.ReadFile("./id_rsa")
	if err != nil {
		panic(err)
	}
	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		panic(err)
	}

	return key
}
