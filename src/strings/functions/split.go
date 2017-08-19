// func Split(S string, sep string) []string

package main

import (
	"fmt"
	"strings"
)

func main() {
	strSlice := strings.Split("a,b,c", ",")
	fmt.Println(strSlice, "\n")

	strSlice = strings.Split("Australia is a country and continent surrounded by the Indian and Pacific oceans.", " ")
	for _, v := range strSlice {
		fmt.Println(v)
	}

	strSlice = strings.Split("abacadaeaf", "a")
	fmt.Println("\n", strSlice)

	strSlice = strings.Split("abacadaeaf", "A")
	fmt.Println("\n", strSlice)

	strSlice = strings.Split("123023403450456056706780789", "0")
	fmt.Println("\n", strSlice)

	strSlice = strings.Split("123023403450456056706780789", ",")
	fmt.Println("\n", strSlice)
}
