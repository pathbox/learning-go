package email

import (
	"encoding/base64"
	"fmt"
	"net/mail"
	"net/smtp"
)

func SendToEvernote(user, password, host, to, subject, body string) error {
	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	from := mail.Address{user, user}
	toMail := mail.Address{to, to}
	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = toMail.String()
	header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", b64.EncodeToString([]byte(subject))) // 这一步，用base64对subject进行编码
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=UTF-8"
	header["Content-Transfer-Encoding"] = "base64"
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + b64.EncodeToString([]byte(body))
	auth := smtp.PlainAuth("", user, password, host)
	err := smtp.SendMail(
		host+":25",
		auth,
		user,
		[]string{toMail.Address},
		[]byte(message),
	)
	if err != nil {
		panic(err)
	}
	return err
}
