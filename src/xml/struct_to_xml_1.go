package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type StringResources struct {
	XMLName        xml.Name         `xml:"resources"`
	ResourceString []ResourceString `xml:"string"`
}

type ResourceString struct {
	XMLName    xml.Name `xml:"string"`
	StringName string   `xml:"name,attr"`
	InnerText  string   `xml:",innerxml"`
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

	for i, item := range result.ResourceString {
		log.Println(item.StringName + "===" + item.InnerText)

		//修改ApplicationName节点的内部文本innerText
		if strings.EqualFold(item.StringName, "ApplicationName") {
			fmt.Println("change innerText")

			//注意修改的不是line对象，而是直接使用result中的真实对象 item 只是一个副本
			result.ResourceString[i].InnerText = "这是新的ApplicationName"
		}
	}

	// 生成 xml
	xmlOut, outErr := xml.MarshalIndent(result, "", "")
	if outErr == nil {
		// 加入XML头
		headerBytes := []byte(xml.Header)
		// 拼接XML头和实际XML内容
		xmlOutData := append(headerBytes, xmlOut...)
		// 写文件
		ioutil.WriteFile("create_xml.xml", xmlOutData, 0775)

		fmt.Println("OK~")
	} else {
		fmt.Println(outErr)
	}

}
