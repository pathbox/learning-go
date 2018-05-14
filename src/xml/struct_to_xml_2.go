package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type Address struct {
	City, State string
}

type Person struct {
	XMLName   xml.Name `xml:"person"` // 最外层的标签
	Id        int      `xml:"id,attr"`
	FirstName string   `xml:"name>first"` // name 标签中内嵌first标签
	LastName  string   `xml:"name>last"`
	Age       int      `xml:"age"` // age 标签， int表示内文本类型为整数 实际用string 对生成xml没有区别
	Height    float32  `xml:"height,omitempty"`
	Married   bool
	Address
	Comment string `xml:",comment"` // 自定义注释
}

func main() {
	v := &Person{Id: 13, FirstName: "John", LastName: "Doe", Age: 42}
	v.Comment = " Need more details. "
	v.Address = Address{"Hanga Roa", "Easter Island"}

	output, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	// os.Stdout.Write(output)
	ioutil.WriteFile("create_xml_2.xml", output, 0755)
}
