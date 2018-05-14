package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

func main() {
	type Address struct {
		City, State string
	}
	type Person struct {
		XMLName   xml.Name `xml:"person"`
		Id        int      `xml:"id,attr"`
		FirstName string   `xml:"name>first"`
		LastName  string   `xml:"name>last"`
		Age       int      `xml:"age"`
		Height    float32  `xml:"height,omitempty"`
		Married   bool
		Address
		Comment string `xml:",comment"`
	}

	v := &Person{Id: 13, FirstName: "John", LastName: "Doe", Age: 42}
	v.Comment = " Need more details. "
	v.Address = Address{"Hanga Roa", "Easter Island"}

	file, _ := os.OpenFile("create_xml_3.xml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)

	// enc := xml.NewEncoder(os.Stdout)
	enc := xml.NewEncoder(file)

	enc.Indent("", "  ")
	if err := enc.Encode(v); err != nil { // Encode 操作的时候，v 会输出到os.Stdout 或者一个文件
		fmt.Printf("error: %v\n", err)
	}
}
