package arena

import (
	"fmt"
	// "strings"
	"testing"
)

func TestArena(t *testing.T) {
	a := NewArena(100)

	for i := 50; i < 100; i++ {
		a.buf[i] = 2
	}

	b1 := a.Make(50)

	b2 := a.Make(30)

	b3 := a.Make(40)

	fmt.Printf("%p %d\n", b1, b1[49])
	fmt.Printf("%p %d\n", b2, b2[29])
	fmt.Printf("%p %d\n", b3, b3[39])
	fmt.Println(b1, b2, b3)
}
