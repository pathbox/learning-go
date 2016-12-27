package main

import (
	"context"
	"fmt"
)

type otherContext struct {
	context.Context
}

type Key struct {
	name string
	age  int
}

func main() {
	sKey := Key{"Curry", 28}
	c0 := context.Background()
	fmt.Println(c0)

	c1 := context.WithValue(c0, "key1", "value1")
	fmt.Printf("key: key1 value: %s\n", c1.Value("key1"))

	c2 := context.WithValue(c1, "k2", sKey)
	switch t := c2.Value("k2").(type) {
	case Key:
		fmt.Printf("key: k2 value:Key{name:%s,age:%d}\n", t.name, t.age)
	default:
		fmt.Println("unknow")
	}
	fmt.Println(c2.Value("k1"))
}
