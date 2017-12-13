package main

import (
	"fmt"

	"github.com/benmanns/goworker"
)

func init() {
	goworker.Register("Hello", helloworker)
}

func helloworker(queue string, args ...interface{}) error {
	fmt.Println("Hello world")
	return nil
}
