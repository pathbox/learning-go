//func HasPrefix(s, prefix string) bool

package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.HasPrefix("Australia", "Aus"))
	fmt.Println(strings.HasPrefix("Australia", "aus"))
	fmt.Println(strings.HasPrefix("Australia", "Jap"))
	fmt.Println(strings.HasPrefix("Australia", ""))
}
