package main

import (
	"fmt"
	"strings"
)

func main() {
	var b strings.Builder // It is a io.Writer
	b.Grow(50)            //Grow grows b's capacity, if necessary, to guarantee space for another n bytes. After Grow(n), at least n bytes can be written to b without another allocation
	for i := 0; i < 3; i++ {
		fmt.Fprintf(&b, "%d...", i)
	}
	b.WriteString("Hello World")

	b.Write([]byte(` Morning`))
	fmt.Println("Len: ", b.Len())
	fmt.Println(b.String())
	b.Reset() // To empty
	fmt.Println(b.String())
}
