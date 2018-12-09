package main

import (
	"fmt"

	"./example1"
)

func main() {
	mycontent := `
    <html>
    <body>
    <h3>
    "Test send to email"
    </h3>
    </body>
    </html>
    `

	email := sendemail.NewEmail("test@163.com",
		"test golang email", mycontent)

	err := sendemail.SendEmail(email)

	fmt.Println(err)
	fmt.Println("Send Done")

}
