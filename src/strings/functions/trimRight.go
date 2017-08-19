// func TrimRight(s string, cutset string) string

package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.TrimRight("0120 2510", "0"))
	fmt.Println(strings.TrimRight("abcd axyz", "a"))
	fmt.Println(strings.TrimRight("abcd axyz", "A"))
	fmt.Println(strings.TrimRight("! Abcd dcbA !", "A"))
	fmt.Println(strings.TrimRight("! Abcd dcbA !", "!"))
	fmt.Println(strings.TrimRight(" Abcd dcbA ", " "))
}
