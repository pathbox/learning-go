// func Trim(s string, cutset string) string

package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.Trim("0120 2510", "0"))
	fmt.Println(strings.Trim("abcd axyz", "a"))
	fmt.Println(strings.Trim("abcd axyz", "A"))
	fmt.Println(strings.Trim("! Abcd dcbA !", "A"))
	fmt.Println(strings.Trim("! Abcd dcbA !", "!"))
	fmt.Println(strings.Trim(" Abcd dcbA ", " "))
}
