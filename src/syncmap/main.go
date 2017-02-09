package main

import (
	"fmt"

	"github.com/DeanThompson/syncmap"
)

func main() {
	m := syncmap.New()
	var value map[string]string
	value = map[string]string{"nice": "hello"}
	m.Set("one", value)
	v, ok := m.Get("one")
	fmt.Println(v, ok) // 1, true

	v, ok = m.Get("not_exist")
	fmt.Println(v, ok) // nil, false

	m.Set("two", 2)
	m.Set("three", "three")

	for item := range m.IterItems() {
		fmt.Println("key:", item.Key, "value:", item.Value)
	}
}
