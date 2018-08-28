package main

import (
	"fmt"

	skiplist "github.com/huandu/skiplist"
)

func main() {
	list := skiplist.New(skiplit.Int) // key's type is int

	// Add some values. Value can be anything.
	list.Set(12, "hello world")
	list.Set(34, 56)

	elem := list.Get(34) // value is stored in elem.Value
	fmt.Println(elem.Value)
	next := elem.Next() // Get next element.
	fmt.Println(next.Value)

	// Or get value directly just like a map
	val, ok := list.GetValue(34)
	fmt.Println(val, ok)

	// Remove an element by i

}
