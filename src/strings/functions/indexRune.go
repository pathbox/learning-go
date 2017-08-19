// func IndexRune(s string, r rune) int

package main

import (
	"fmt"
	"strings"
)

func main() {
	var s, t, u rune
	t = 'l'
	fmt.Println(strings.IndexRune("australia", t))
	fmt.Println(strings.IndexRune("LONDON", t))
	fmt.Println(strings.IndexRune("JAPAN", t))

	s = 1
	fmt.Println(strings.IndexRune("5221-JAPAN", s))

	u = '1'
	fmt.Println(strings.IndexRune("5221-JAPAN", u))
}
