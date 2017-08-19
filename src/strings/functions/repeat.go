// func Repeat(s string, count int) string

package main

import (
	"fmt"
	"strings"
)

func main() {
	textString := "Japan"
	repString := strings.Repeat(textString, 5)
	fmt.Println(repString)

	textString = " A " // char with space on both side
	repString = strings.Repeat(textString, 5)
	fmt.Println(repString) // Repeat space also

	fmt.Println("ba" + strings.Repeat("na", 2))
	fmt.Println("111" + strings.Repeat("22", 2))
	fmt.Println("111" + strings.Repeat(" ", 2))
	fmt.Println(strings.Repeat("=", 50))
}
