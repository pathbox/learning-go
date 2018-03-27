package main

import (
	"fmt"
	"regexp"
)

func main() {
	b := []byte("data[title]")
	pat := `^data`
	reg1 := regexp.MustCompile(pat)
	r := reg1.Find(b)
	fmt.Println(string(r))
	rb := reg1.Match(b)
	fmt.Println(rb)
}
