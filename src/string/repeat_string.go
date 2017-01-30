package main

import (
	"fmt"
	"strings"
)

func main() {
	var origS string = "Hi here! "
	var newS string

	newS = strings.Repeat(origS, 3)
	fmt.Println(newS)
}
