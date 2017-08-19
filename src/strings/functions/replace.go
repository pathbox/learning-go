// func Replace(s, old, new string, n int) string

package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.Replace("Australia Japan Canada Indiana", "an", "anese", 0))
	fmt.Println(strings.Replace("Australia Japan Canada Indiana", "an", "anese", 1))
	fmt.Println(strings.Replace("Australia Japan Canada Indiana", "an", "anese", 2))
	fmt.Println(strings.Replace("Australia Japan Canada Indiana", "an", "anese", -1))
	fmt.Println(strings.Replace("1234534232132", "1", "0", -1))
	// case-sensitive
	fmt.Println(strings.Replace("Australia Japan Canada Indiana", "AN", "anese", -1))
}
