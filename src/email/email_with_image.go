package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"strings"
	"time"
)

type SendMail struct {
	user     string
	password string
	host     string
	port     string
	auth     smtp.Auth
}

type Attachment struct {
	name        []string
	contentType string
	withFile    bool
}

type Message struct {
	from        string
	to          []string
	cc          []string
	bcc         []string
	subject     string
	body        string
	contentType string
	attachment  Attachment
}

func (mail *SendMail) Auth() {
	mail.auth = smtp.PlainAuth("", mail.user, mail.password, mail.host)
}

func (mail SendMail) Send(message Message) error {
	mail.Auth()
	buffer := bytes.NewBuffer(nil)
	boundary := "GoBoundary"
	Header := make(map[string]string)
	Header["From"] = message.from
	Header["To"] = strings.Join(message.to, ";")
	Header["Cc"] = strings.Join(message.cc, ";")
	Header["Bcc"] = strings.Join(message.bcc, ";")
	Header["Subject"] = message.subject
	Header["Content-Type"] = "multipart/related;boundary=" + boundary
	Header["Date"] = time.Now().String()
	mail.writeHeader(buffer, Header)

	var imgsrc string
	if message.attachment.withFile {
		//多图片发送
		for _, graphname := range message.attachment.name {
			attachment := "\r\n--" + boundary + "\r\n"
			attachment += "Content-Transfer-Encoding:base64\r\n"
			attachment += "Content-Type:" + message.attachment.contentType + ";name=\"" + graphname + "\"\r\n"
			attachment += "Content-ID: <" + graphname + "> \r\n\r\n"
			buffer.WriteString(attachment)

			//拼接成html
			imgsrc += "<p><img src=\"cid:" + graphname + "\" height=500 width=500></p><br>\r\n\t\t\t"

			defer func() {
				if err := recover(); err != nil {
					fmt.Printf(err.(string))
				}
			}()
			mail.writeFile(buffer, graphname)
		}
	}

	//需要在正文中显示的html格式
	var template = `
    <html>
        <body>
            <p>text:%s</p><br>
            %s
        </body>
    </html>
    `
	var content = fmt.Sprintf(template, message.body, imgsrc)
	body := "\r\n--" + boundary + "\r\n"
	body += "Content-Type: text/html; charset=UTF-8 \r\n"
	body += content
	buffer.WriteString(body)

	buffer.WriteString("\r\n--" + boundary + "--")
	fmt.Println(buffer.String())
	smtp.SendMail(mail.host+":"+mail.port, mail.auth, message.from, message.to, buffer.Bytes())
	return nil
}

func (mail SendMail) writeHeader(buffer *bytes.Buffer, Header map[string]string) string {
	header := ""
	for key, value := range Header {
		header += key + ":" + value + "\r\n"
	}
	header += "\r\n"
	buffer.WriteString(header)
	return header
}

func (mail SendMail) writeFile(buffer *bytes.Buffer, fileName string) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	payload := make([]byte, base64.StdEncoding.EncodedLen(len(file)))
	base64.StdEncoding.Encode(payload, file)
	buffer.WriteString("\r\n")
	for index, line := 0, len(payload); index < line; index++ {
		buffer.WriteByte(payload[index])
		if (index+1)%76 == 0 {
			buffer.WriteString("\r\n")
		}
	}
}

func main() {
	mail := &SendMail{user: "test@163.com", password: "test", host: "smtp.163.com", port: "25"}
	message := Message{
		from:        "test@163.com",
		to:          []string{"test@163.com"},
		cc:          []string{},
		bcc:         []string{},
		subject:     "Test Email with golang", //邮件标题
		body:        "msg body!",              //正文内容
		contentType: "text/plain;charset=utf-8",
		attachment: Attachment{
			name:        []string{"/path/to/image.png"}, //可以放入多张图片
			contentType: "image/png",
			withFile:    true,
		},
	}
	mail.Send(message)
}
