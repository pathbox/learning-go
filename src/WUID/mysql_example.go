package main

import (
	"fmt"

	"github.com/edwingeng/wuid/mysql"
)

func main() {
	g := wuid.NewWUID("default", nil)
	g.LoadH24FromMysql("127.0.0.1:3306", "root", "", "test", "wuid")

	for i := 0; i < 1000000; i++ {
		fmt.Printf("%#016x\n", g.Next())
		fmt.Println(g.Next())
	}
}
