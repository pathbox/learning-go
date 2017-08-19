// func FieldsFunc(s string, f func(rune) bool) []string

package main

import (
	"fmt"
	"strings"
	"unicode"
)

func main() {

	x := func(c rune) bool {
		return !unicode.IsLetter(c)
	}
	strArray := strings.FieldsFunc(`Australia major cities – Sydney, Brisbane,
                                 Melbourne, Perth, Adelaide – are coastal`, x)
	for _, v := range strArray {
		fmt.Println(v)
	}

	fmt.Println("\n*****************Split by number*******************\n")

	y := func(c rune) bool {
		return unicode.IsNumber(c)
	}
	testArray := strings.FieldsFunc(`1 Sydney Opera House.2 Great Barrier Reef.3 Uluru-Kata Tjuta National Park.4 Sydney Harbour Bridge.5 Blue Mountains National Park.6 Melbourne.7 Bondi Beach`, y)
	for _, w := range testArray {
		fmt.Println(w)
	}

}

// Australia
// major
// cities
// Sydney
// Brisbane
// Melbourne
// Perth
// Adelaide
// are
// coastal

// *****************Split by number*******************

//  Sydney Opera House.
//  Great Barrier Reef.
//  Uluru-Kata Tjuta National Park.
//  Sydney Harbour Bridge.
//  Blue Mountains National Park.
//  Melbourne.
//  Bondi Beach
