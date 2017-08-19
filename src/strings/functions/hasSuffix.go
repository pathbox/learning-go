// func HasSuffix(s, prefix string) bool

package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.HasSuffix("Australia", "lia"))
	fmt.Println(strings.HasSuffix("Australia", "A"))
	fmt.Println(strings.HasSuffix("Australia", "LIA"))
	fmt.Println(strings.HasSuffix("123456", "456"))
	fmt.Println(strings.HasSuffix("Australia", ""))
}
