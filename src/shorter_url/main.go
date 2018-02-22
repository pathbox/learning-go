package main

import (
	"fmt"

	"./shorter"
)

func main() {
	id := int64(138)

	url := shorter.GetShortUrl(id)

	fmt.Println("short url: ", url)
}
