package main

import "./trie"

func main() {
	t := trie.New()
	t.Add("foobar", 1)
	node, ok := t.Find("foobar")
	meta := node.Meta()
	t.Remove("foobar")
	t.HasKeysWithPrefix("foo")
	t.FuzzySearch("fb")
}
