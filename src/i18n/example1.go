package main

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func init() {
	message.SetString(language.Chinese, "%s want to %s.", "%s 想要去 %s.")
	message.SetString(language.AmericanEnglish, "%s went to %s.", "%s is in %s.")
	message.SetString(language.Chinese, "%s has been stolen.", "%s 被偷走了.")
	message.SetString(language.AmericanEnglish, "%s has been stolen.", "%s has been stolen.")
	message.SetString(language.Chinese, "How are you?", "最近怎么样?")
}

func main() {
	p := message.NewPrinter(language.Chinese)
	p.Printf("%s want to %s.", "Cary", "日本旅游")
	fmt.Println()
	p.Printf("%s has been stolen.", "宝藏")
	fmt.Println()

	p = message.NewPrinter(language.AmericanEnglish)
	p.Printf("%s went to %s.", "Peter", "England")
	fmt.Println()
	p.Printf("%s has been stolen.", "The Gem")
}
