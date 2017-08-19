// func Count(s, sep string) int

package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.Count("Australia", "a"))
	fmt.Println(strings.Count("Australia", "A"))
	fmt.Println(strings.Count("Australia", "M"))
	fmt.Println(strings.Count("Japanese", "Japan")) // 1
	fmt.Println(strings.Count("Japan", "Japanese")) // 0
	fmt.Println(strings.Count("Shell-25152", "-25"))
	fmt.Println(strings.Count("Shell-25152", "-21"))
	fmt.Println(strings.Count("test", "")) // length of string + 1
	fmt.Println(strings.Count("test", " "))
}
