// 这两个方法是当编解码中有一个字段是interface{}的时候需要对interface{}的可能产生的类型进行注册

package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

type P struct {
	X, Y, Z int
	Name    interface{}
}

type Q struct {
	X, Y *int32
	Name interface{}
}

type Inner struct {
	Test int
}

func main() {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network) // Will write to network.
	dec := gob.NewDecoder(&network) // Will read from network.

	gob.Register(Inner{})

	inner := Inner{100}

	err := enc.Encode(P{1, 2, 3, inner})
	if err != nil {
		log.Fatal("encode error:", err)
	}

	var q Q
	err = dec.Decode(&q)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	fmt.Println(q)
	fmt.Printf("%v: {%d,%d}\n", q.Name, *q.X, *q.Y)

}
