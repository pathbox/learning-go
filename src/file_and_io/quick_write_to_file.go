package main

import (
	"io/ioutil"
	"log"
)

func main() {
	err := ioutil.WriteFile("files/empty_new.txt", []byte("Hello Kitty"), 0666) // 这是覆盖的写入，会覆盖之前文件中的内容
	if err != nil {
		log.Fatalln(err)
	}
}
