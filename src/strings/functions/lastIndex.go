// func LastIndex(s, sep string) int

package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.LastIndex("japanese", "a")) // position j=0,a=1,p=2,a=3
	fmt.Println(strings.LastIndex("japanese", "A"))
	fmt.Println(strings.LastIndex("Australia Australian", "lia"))
	fmt.Println(strings.LastIndex("AUSTRALIA AUSTRALIAN", "lia"))
	fmt.Println(strings.LastIndex("1234567890 1234567890", "0"))
	fmt.Println(strings.LastIndex("1234567890 1234567890", "00"))
	fmt.Println(strings.LastIndex("1234567890 1234567890", "123"))
}
