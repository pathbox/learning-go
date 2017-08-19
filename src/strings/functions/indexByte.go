// func IndexByte(s string, c byte) int

package main

import (
	"fmt"
	"strings"
)

func main() {
	var s, t, u byte
	t = 'l'
	fmt.Println(strings.IndexByte("australia", t))
	fmt.Println(strings.IndexByte("LONDON", t))
	fmt.Println(strings.IndexByte("JAPAN", t))

	s = 1
	fmt.Println(strings.IndexByte("5221-JAPAN", s))

	u = '1'
	fmt.Println(strings.IndexByte("5221-JAPAN", u))
}
