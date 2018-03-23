package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/mdp/qrterminal"
)

func main() {
	// qrterminal.Generate("https://github.com/mdp/qrterminal", qrterminal.L, os.Stdout)
	// qrterminal.Generate("https://github.com/mdp/qrterminal", qrterminal.M, os.Stdout)
	// qrterminal.Generate("https://github.com/mdp/qrterminal", qrterminal.H, os.Stdout)

	buf := &bytes.Buffer{}

	path := "/home/user/zq"
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err)
	}
	qrterminal.Generate("https://github.com/mdp/qrterminal", qrterminal.H, buf)

	fmt.Println(len(buf.Bytes()))
	fmt.Println(string(buf.Bytes()))
	n, err := f.Write(buf.Bytes())

	fmt.Println("n", n)

	if err != nil {
		fmt.Println(err)
	}
}
