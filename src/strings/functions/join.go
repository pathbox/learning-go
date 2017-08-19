// func Join(stringSlice []string, sep string) string

package main

import (
	"fmt"
	"strings"
)

func main() {
	// Slice of strings
	textString := []string{"Australia", "Japan", "Canada"}
	fmt.Println(strings.Join(textString, "-"))

	// Slice of strings
	textNum := []string{"1", "2", "3", "4", "5"}
	fmt.Println(strings.Join(textNum, ""))
}
