package main

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"bufio"
	"bytes"
	"fmt"
	"image/png"
	"io/ioutil"
	"os"
)

func display(key *otp.Key, data []byte) {
	fmt.Printf("Issuer:       %s\n", key.Issuer())
	fmt.Printf("Account Name: %s\n", key.AccountName())
	fmt.Printf("Secret:       %s\n", key.Secret())
	fmt.Println("Writing PNG to qr-code.png....")
	ioutil.WriteFile("qr-code.png", data, 0644)
	fmt.Println("")
	fmt.Println("Please add your TOTP to your OTP Application now!")
	fmt.Println("")
}

func promptForPasscode() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter password")
	text, _ := reader.ReadString('\n')
	return text
}

func main() {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "example.com",
		AccountName: "joe@example.com",
	})
	if err != nil {
		panic(err)
	}

	// convert TOTP key into a PNG
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		panic(err)
	}
	png.Encode(&buf, img)

	// display the QR code to the user
	display(key, buf.Bytes())

	// Now Validata that the user's successfully added the passcode.

	fmt.Println("Validating TOTP......")
	passcode := promptForPasscode()
	fmt.Println("passcode:", passcode)
	fmt.Println("key:", key.Secret())
	valid := totp.Validate(passcode, key.Secret())
	if valid {
		println("Valid passcode!")
	} else {
		println("Invalid passcode!")
		os.Exit(1)
	}
}
