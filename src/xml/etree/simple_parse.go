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
		t := node.SelectAttr("type")
		id := node.SelectAttr("id")
		fmt.Println(t.Value)
		fmt.Println(id.Value)
		e := node.SelectElement("on_succ")
		if e != nil {

			// fmt.Println(e.Tag)
		}
	}
}
