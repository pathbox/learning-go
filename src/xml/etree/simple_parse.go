package main

import (
	"fmt"

	"github.com/beevik/etree"
)

func main() {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile("simple.xml"); err != nil {
		panic(err)
	}

	root := doc.SelectElement("udesk_ivr")
	fmt.Println("Root element", root.Tag)
	for _, node := range root.SelectElements("node") {
		// for _, attr := range node.Attr {
		// 	fmt.Println(attr.Key, attr.Value, attr.Space)
		// }
		// t := node.SelectAttr("type")
		// id := node.SelectAttr("id")
		// fmt.Println(t.Value)
		// fmt.Println(id.Value)
		e := node.SelectElement("on_succ")
		if e != nil {
			// fmt.Println(e.Tag)
			actions := e.SelectElements("action")
			for _, action := range actions {
				gotoID := action.SelectAttrValue("data", "")
				fmt.Println(gotoID)
			}
		}
		for _, el := range node.ChildElements() {
			fmt.Println(el.Tag)
		}

		attrs := node.SelectElement("attrs")
		if attrs != nil {
			fmt.Println("===========", attrs.SelectAttrValue("type", ""))
			fmt.Println("===========", attrs.SelectAttrValue("content", ""))
		}
	}
}
