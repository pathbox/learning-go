package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
)

type StringResources struct { // 嵌套的 xml结构
	XMLName        xml.Name         `xml:"resources"` // xml标签名称
	ResourceString []ResourceString `xml:"string"`
}

type ResourceString struct {
	XMLName    xml.Name `xml:"string"`    // xml标签名称
	StringName string   `xml:"name,attr"` // xml 标签中的属性
	InnerText  string   `xml:",innerxml"` // 标签直接的文本
}

func main() {
	XMLContent, err := ioutil.ReadFile("example1.xml")
	if err != nil {
		log.Fatal(err)
	}

	var result StringResources
	err = xml.Unmarshal(XMLContent, &result) // 反序列化到struct
	if err != nil {
		log.Fatal(err)
	}

	log.Println(result)
	log.Println(result.ResourceString)

	for _, item := range result.ResourceString {
		log.Println(item.StringName + "===" + item.InnerText)
	}
}
