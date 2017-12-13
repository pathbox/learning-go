package main

import (
	"fmt"

	"github.com/benmanns/goworker"
)

func main() {
	if err := goworker.Work(); err != nil {
		fmt.Println("Error", err)
	}
}
