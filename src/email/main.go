package main

import (
    "fmt"
    "./example1"
)

func main() {
    mycontent := " my dear"

    email := sendemail.NewEmail("luodaokai@udesk.cn",
        "test golang email", mycontent)

    err := sendemail.SendEmail(email)

    fmt.Println(err)

}