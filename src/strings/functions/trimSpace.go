// func TrimSpace(s string) string

package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.TrimSpace(" New Zealand is a country in the southwestern Pacific Ocean "))
	fmt.Println(strings.TrimSpace(" \t\n  New Zealand is a country in the southwestern Pacific Ocean \t\n "))
	fmt.Println(strings.TrimSpace(" \t\n\r\x0BNew Zealand is a country in the southwestern Pacific Ocean\t\n "))
}
