// func SplitN(s, sep string, n int) []string

package main

import (
	"fmt"
	"strings"
)

func main() {
	strSlice := strings.SplitN("a,b,c", ",", 0)
	fmt.Println(strSlice, "\n")

	strSlice = strings.SplitN("a,b,c", ",", 1)
	fmt.Println(strSlice, "\n")

	strSlice = strings.SplitN("a,b,c", ",", 2)
	fmt.Println(strSlice, "\n")

	strSlice = strings.SplitN("a,b,c", ",", 3)
	fmt.Println(strSlice, "\n")

	strSlice = strings.SplitN("Australia is a country and continent surrounded by the Indian and Pacific oceans.", " ", -1)
	for _, v := range strSlice {
		fmt.Println(v)
	}
	strSlice = strings.SplitN("123023403450456056706780789", "0", 5)
	fmt.Println("\n", strSlice)
}
