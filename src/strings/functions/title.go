// func Title(s string) string

package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.Title("Germany is a Western European country with a landscape of forests, rivers, mountain ranges and North Sea beaches."))
	fmt.Println(strings.Title("BAVARIA HESSE BRANDENBURG SAARLAND"))
	fmt.Println(strings.Title("towns and cities"))
}
