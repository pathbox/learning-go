// func ToTitle(s string) string

package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.ToTitle("Germany is a Western European country with a landscape of forests, rivers, mountain ranges and North Sea beaches."))
	fmt.Println(strings.ToTitle("BAVARIA HESSE BRANDENBURG SAARLAND"))
	fmt.Println(strings.ToTitle("towns and cities"))
}
