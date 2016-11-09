package main

import (
	"fmt"
)

func main() {
	x := 100
	switch {
	case x >= 0:
		if x == 0 {
			break
		}
		fmt.Println(x)
	}
}
