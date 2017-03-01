package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	key1 := "Hello" + ":" + "World"
	key2 := fmt.Sprintf("%s:%s", "Hello", "World")
	key3 := strings.Join([]string{"Hello", "World"}, ":")

	l := len("Hello") + len("World") + 2
	buf := make([]byte, 0, l)
	w := bytes.NewBuffer(buf)
	w.WriteString("Hello")
	w.WriteRune(':')
	w.WriteString("World")
	key4 := w.String()

	fmt.Println(key1)
	fmt.Println(key2)
	fmt.Println(key3)
	fmt.Println(key4)
}
