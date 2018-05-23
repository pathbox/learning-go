package main

import (
	"fmt"

	"github.com/rs/xid"
)

func main() {
	guid1 := xid.New()

	fmt.Println(guid1.String())
	fmt.Println(string(guid1.Machine()))
	fmt.Println(guid1.Pid())
	fmt.Println(guid1.Time())
	fmt.Println(guid1.Counter())

	guid2 := xid.New()

	fmt.Println(guid2.String())
	fmt.Println(string(guid1.Machine()))
	fmt.Println(guid2.Pid())
	fmt.Println(guid2.Time())
	fmt.Println(guid2.Counter())

}
