package main

import (
	"fmt"

	"github.com/beevik/etree"
)

func main() {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile("/Users/pathbox/code/learning-go/src/xml/etree/bookstore.xml"); err != nil {
		panic(err)
	}

	root := doc.SelectElement("bookstore")
	fmt.Println("ROOT element:", root.Tag)

	for _, book := range root.SelectElements("book") {
		fmt.Println("Child element:", book.Tag)
		if title := book.SelectElement("title"); title != nil {
			lang := title.SelectAttrValue("lang", "unknown")
			fmt.Printf("  TITLE: %s (%s)\n", title.Text(), lang)
		}
		for _, attr := range book.Attr {
			fmt.Printf("  ATTR: %s=%s\n", attr.Key, attr.Value)
		}
	}

	for _, t := range doc.FindElements("//book[@category='WEB']/title") {
		fmt.Println("Title: ", t.Text())
	}
}
