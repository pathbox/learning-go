// func TrimLeft(s string, cutset string) string

package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.TrimLeft("0120 2510", "0"))
	fmt.Println(strings.TrimLeft("abcd axyz", "a"))
	fmt.Println(strings.TrimLeft("abcd axyz", "A"))
	fmt.Println(strings.TrimLeft("! Abcd dcbA !", "A"))
	fmt.Println(strings.TrimLeft("! Abcd dcbA !", "!"))
	fmt.Println(strings.TrimLeft(" Abcd dcbA ", " "))
}
