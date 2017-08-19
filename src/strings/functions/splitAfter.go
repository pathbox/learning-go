// func SplitAfter(S String, sep string) []string

package main

import (
	"fmt"
	"strings"
)

func main() {
	strSlice := strings.SplitAfter("a,b,c", ",")
	fmt.Println(strSlice, "\n")

	strSlice = strings.SplitAfter("Australia is a country and continent surrounded by the Indian and Pacific oceans.", " ")
	for _, v := range strSlice {
		fmt.Println(v)
	}

	strSlice = strings.SplitAfter("abacadaeaf", "a")
	fmt.Println("\n", strSlice)

	strSlice = strings.SplitAfter("abacadaeaf", "A")
	fmt.Println("\n", strSlice)

	strSlice = strings.SplitAfter("123023403450456056706780789", "0")
	fmt.Println("\n", strSlice)

	strSlice = strings.SplitAfter("123023403450456056706780789", ",")
	fmt.Println("\n", strSlice)
}
