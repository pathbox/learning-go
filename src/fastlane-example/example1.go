package main

import (
	"fmt"

	"github.com/tidwall/fastlane"
)

func main() {
	var ch fastlane.Chan

	go func() {
		for i := 0; i < 10; i++ {
			ch.Send(i)
		}

	}()

	for i := 0; i < 10; i++ {
		v := ch.Recv()
		fmt.Println(v.(int))
	}

}
