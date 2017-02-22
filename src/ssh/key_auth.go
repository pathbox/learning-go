package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

func main() {
	fmt.Println("Here ssh auth start")
	file := "/home/user/.ssh/id_rsa"

	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		panic(err)
	}
	authMethod := ssh.PublicKeys(key)

	config := &ssh.ClientConfig{ // one: config config
		User: "user",
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		Config: ssh.Config{
			Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "arcfour256",
				"arcfour128", "aes128-cbc", "3des-cbc", "blowfish-cbc", "cast128-cbc",
				"aes192-cbc", "aes256-cbc", "arcfour"},
		},
	}

	client, err := ssh.Dial("tcp", "server_address_or_ip", config) // two: Dial conn to create client
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

// func PublicKeyFile(file string) ssh.AuthMethod {
// 	buffer, err := ioutil.ReadFile(file)
// 	if err != nil {
// 		return nil
// 	}

// 	key, err := ssh.ParsePrivateKey(buffer)
// 	if err != nil {
// 		return nil
// 	}
// 	return ssh.PublicKeys(key)
// }

// func SSHAgent() ssh.AuthMethod {
//   if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
//     return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
//   }
//   return nil
// }
// Then you can use the function to instanciate the client config in the following way:

// sshConfig := &ssh.ClientConfig{
//   User: "your_user_name",
//   Auth: []ssh.AuthMethod{
//     SSHAgent()
//   },
// }
// Note that you can add your certificate to the SSH agent by using the following command:

// $ ssh-add /path/to/your/private/certificate/file
