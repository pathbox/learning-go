package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Person struct {
	Id   int
	Name string
	Age  int
}

func (p *Person) Write(w io.Writer) {
	b, _ := json.Marshal(*p)
	w.Write(b)
}

func main() {
	p := &Person{Id: 1, Name: "Joe", Age: 27}
	var b bytes.Buffer

	p.Write(&b)

	fmt.Println(b.String())

}
