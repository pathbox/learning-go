package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const youDaoUrl string = "http://dict.youdao.com/fsearch?q="

// 构造三个struct，并且直接有关联
type YodaoDict struct {
	ReturnPhrase       string              `xml:"return-phrase"`
	CustomTranslations []CustomTranslation `xml:"custom-translation"`
}

type CustomTranslation struct {
	Title        string        "custom-dict"
	Translations []Translation `xml:"translation"`
}

type Translation struct {
	Content string `xml:"content"`
}

func main() {
	args := os.Args

	if args == nil || len(args) != 2 {
		fmt.Println("USAGE: youdao [word] ...")
		os.Exit(1)
	}

	res, err := http.Get(youDaoUrl + args[1])
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	res.Body.Close()

	content := YodaoDict{}
	err = xml.Unmarshal(data, &content)
	if err != nil {
		log.Fatal(err)
	}

	customTrans := content.CustomTranslations
	word := content.ReturnPhrase
	if len(customTrans) == 0 {
		log.Fatal("ERROR: no content from Youdao")
	}
	fmt.Printf("%s:\n", word)

	trans := customTrans[0].Translations
	if len(trans) == 0 {
		log.Fatal("ERROR: no content from Youdao")
	}

	for i := 0; i < len(trans); i++ {
		fmt.Printf("%s\n", trans[i].Content)
	}
}
