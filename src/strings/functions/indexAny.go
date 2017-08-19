// func IndexAny(s, chars string) int

package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.IndexAny("australia", "japan")) // a position
	fmt.Println(strings.IndexAny("japan", "pen"))       // p position
	fmt.Println(strings.IndexAny("mobile", "one"))      // o position
	fmt.Println(strings.IndexAny("123456789", "4"))     // 4 position
	fmt.Println(strings.IndexAny("123456789", "0"))     // 0 position
}
