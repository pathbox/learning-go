package main

import (
	"fmt"
)

func foo(i *int) {
	fmt.Println(*i)
}

func test() (err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("exception occured, err: ", err)
		}
	}()

	for i := 0; i <= 10; i++ {
		foo(&i)
		panic("stop here") // stop here
	}
	foo(nil)
	return
}

func main() {
	test()
}
