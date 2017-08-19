// func TrimPrefix(S string, prefix string) string

package main

import (
	"fmt"
	"strings"
)

func main() {
	var s string
	s = "Australia Canada Japan"
	s = strings.TrimPrefix(s, "Australia")
	s = strings.TrimSpace(s)
	fmt.Println(s)

	s = "Australia Canada Japan"
	s = strings.TrimPrefix(s, "australia")
	fmt.Println(s)

	s = "\nAustralia-Canada-Japan"
	s = strings.TrimPrefix(s, "\n")
	fmt.Println(s)

	s = "\tAustralia-Canada-Japan"
	s = strings.TrimPrefix(s, "\t")
	fmt.Println(s)
}
