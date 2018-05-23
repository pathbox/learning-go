package main

import (
	"fmt"
	"time"
)

type User struct {
	Name string
	Age  int
}

func main() {

	m := make(map[int]*User)

	for i := 0; i < 1000000; i++ {
		user := &User{
			Name: "Cary",
			Age:  i,
		}

		m[i] = user
	}
	fmt.Println("Size Map: ", len(m))

	time.Sleep(5 * time.Second)
	for i := 0; i < 1000000; i++ {
		// m[i] = nil
		delete(m, i)

	}
	fmt.Println("delete done")
	fmt.Println("After delete size: ", len(m))
	// debug.FreeOSMemory()
	fmt.Println("free memory")
	time.Sleep(1000 * time.Second)

}
