// func EqualFold(s, t string) bool

package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.EqualFold("Australia", "AUSTRALIA"))
	fmt.Println(strings.EqualFold("Australia", "aUSTRALIA"))
	fmt.Println(strings.EqualFold("Australia", "Australia"))
	fmt.Println(strings.EqualFold("Australia", "Aus"))
	fmt.Println(strings.EqualFold("Australia", "Australia & Japan"))
	fmt.Println(strings.EqualFold("JAPAN-1254", "japan-1254"))
	fmt.Println(strings.EqualFold(" ", " "))  // single space both side
	fmt.Println(strings.EqualFold(" ", "  ")) // double space right side
}
