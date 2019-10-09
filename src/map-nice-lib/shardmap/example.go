package main

import "github.com/tidwall/shardmap"

func main() {
	var m shardmap.Map
	m.Set("Hello", "Dolly!")
	val, _ := m.Get("Hello")
	fmt.Printf("%v\n", val)
	val, _ = m.Delete("Hello")
	fmt.Printf("%v\n", val)
	val, _ = m.Get("Hello")
	fmt.Printf("%v\n", val)
}