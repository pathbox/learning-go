package main

import (
    "github.com/go-gomail/gomail"
)

func main() {
    m := gomail.NewMessage()
    m.SetAddressHeader("From", "550755606@qq.com", "一去、二三里")  // 发件人
    m.SetHeader("To",  // 收件人
        m.FormatAddress("********@163.com", "乔峰"),
        m.FormatAddress("********@qq.com", "郭靖"),
    )
    m.SetHeader("Subject", "Gomail")  // 主题
    m.SetBody("text/html", "Hello <a href = \"http://blog.csdn.net/liang19890820\">一去丶二三里</a>")  // 正文

    d := gomail.NewPlainDialer("smtp.qq.com", 465, "550755606@qq.com", "*********")  // 发送邮件服务器、端口、发件人账号、发件人密码
    if err := d.DialAndSend(m); err != nil {
        panic(err)
    }
}

// m := gomail.NewMessage()
// m.SetHeader("From", "alex@example.com")
// m.SetHeader("To", "bob@example.com", "cora@example.com")
// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
// m.SetHeader("Subject", "Hello!")
// m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
// m.Attach("/home/Alex/lolcat.jpg")

// d := gomail.NewDialer("smtp.example.com", 587, "user", "123456")

// // Send the email to Bob, Cora and Dan.
// if err := d.DialAndSend(m); err != nil {
//     panic(err)
// }