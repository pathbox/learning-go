package main

import (
	"fmt"

	"gopkg.in/vmihailenco/msgpack.v2"
)

func main() {
	b, err := msgpack.Marshal(true)
	if err != nil {
		panic(err)
	}
	fmt.Println(b)

	var out bool
	err = msgpack.Unmarshal([]byte{0xc3}, &out)
	if err != nil {
		panic(err)
	}

	fmt.Println(out)
	// Output: []byte{0xc3}
	// true
}


