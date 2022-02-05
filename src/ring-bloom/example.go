package main

import (
	"fmt"

	"./ring"
)

func main() {
	// Support up to 100 elements with less than 1% false positives
	r, err := ring.Init(100, 0.01)
	if err != nil {
		// error will only occur if parameters are set incorrectly
		panic(err)
	}

	data := []byte("hello")

	// check if data is in ring
	fmt.Printf("%s in ring :: %t\n", data, r.Test(data))

	// add data to ring
	r.Add(data)
	fmt.Printf("%s in ring :: %t\n", data, r.Test(data))

	// reset ring
	r.Reset()
	fmt.Printf("%s in ring :: %t\n", data, r.Test(data))
}
