package main

import "github.com/cheynewallace/tabby"

func main() {
	t := tabby.New()
	t.AddHeader("NAME", "TITLE", "DEPARTMENT")
	t.AddLine("John Smith", "Developer", "Engineering")
	t.Print()
}
