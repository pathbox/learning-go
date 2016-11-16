package main

import (
	"fmt"
	"os"

	"github.com/tinylib/msgp/msgp"
)

//go:generate msgp

//msgp:tuple Primitive

type Primitive struct {
	One   int
	Two   uint
	Three float64
}

func main() {
	p := Primitive{1, 2, 3.0}
	out, _ := p.MarshalMsg(nil)
	msgp.UnmarshalAsJSON(os.Stdout, out)
	fmt.Println()
	p1 := new(Primitive)
	left, err := p1.UnmarshalMsg(out)
	fmt.Println(left, err, p1)
}
