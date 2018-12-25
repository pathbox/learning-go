package main

import (
	"log"
	"net/mail"
	"net/smtp"

	"github.com/scorredoira/email"
)

func main() {
	m := email.NewHTMLMessage("This is subject呼啦", "<h1>Nice to meet you早上好</h1>")
	m.From = mail.Address{Name: "From", Address: "test@163.com"}
	m.To = []string{"test@163.com"}

	if err := m.Attach("test.docx"); err != nil {
		log.Fatal(err)
	}
	if err := m.Attach("test.png"); err != nil {
		log.Fatal(err)
	}

	m.AddHeader("X-CUSTOMER-id", "xxxxx")

	auth := smtp.PlainAuth("", "test@163.com", "pwd", "smtp.163.com")
	if err := email.Send("smtp.163.com:25", auth, m); err != nil {
		log.Println("send error")
		log.Fatal(err)
	}

	log.Println("Done~")
}
